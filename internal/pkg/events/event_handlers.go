package events

import (
	"learn/internal/model"
	"learn/internal/pkg/random"
	"learn/internal/repository"
	"log/slog"
)

// OrderPaidEventHandler handles the OrderPaidEvent
type OrderPaidEventHandler struct {
	orderRepo repository.OrderRepository
	logger    *slog.Logger
}

// NewOrderPaidEventHandler creates a new OrderPaidEventHandler
func NewOrderPaidEventHandler(orderRepo repository.OrderRepository, logger *slog.Logger) *OrderPaidEventHandler {
	return &OrderPaidEventHandler{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

// Handle processes the OrderPaidEvent
func (h *OrderPaidEventHandler) Handle(event Event) {
	if orderPaidEvent, ok := event.(OrderPaidEvent); ok {
		order, err := h.orderRepo.GetOrderByID(orderPaidEvent.OrderID)
		if err != nil {
			h.logger.Error("Failed to get order for payment processing",
				slog.Uint64("order_id", uint64(orderPaidEvent.OrderID)),
				slog.String("error", err.Error()))
			return
		}

		// Update order status to paid
		order.Status = model.OrderPaid
		if err := h.orderRepo.UpdateOrder(order); err != nil {
			h.logger.Error("Failed to update order status to paid",
				slog.Uint64("order_id", uint64(orderPaidEvent.OrderID)),
				slog.String("error", err.Error()))
			return
		}

		h.logger.Info("Order status updated to paid",
			slog.Uint64("order_id", uint64(orderPaidEvent.OrderID)),
			slog.Int64("total_price", orderPaidEvent.TotalPrice))
	}
}

// OrderCancelledEventHandler handles the OrderCancelledEvent
type OrderCancelledEventHandler struct {
	orderRepo repository.OrderRepository
	logger    *slog.Logger
}

// NewOrderCancelledEventHandler creates a new OrderCancelledEventHandler
func NewOrderCancelledEventHandler(orderRepo repository.OrderRepository, logger *slog.Logger) *OrderCancelledEventHandler {
	return &OrderCancelledEventHandler{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

// Handle processes the OrderCancelledEvent
func (h *OrderCancelledEventHandler) Handle(event Event) {
	if orderCancelledEvent, ok := event.(OrderCancelledEvent); ok {
		order, err := h.orderRepo.GetOrderByID(orderCancelledEvent.OrderID)
		if err != nil {
			h.logger.Error("Failed to get order for cancellation",
				slog.Uint64("order_id", uint64(orderCancelledEvent.OrderID)),
				slog.String("error", err.Error()))
			return
		}

		// Update order status to cancelled
		order.Status = model.OrderCancelled
		if err := h.orderRepo.UpdateOrder(order); err != nil {
			h.logger.Error("Failed to update order status to cancelled",
				slog.Uint64("order_id", uint64(orderCancelledEvent.OrderID)),
				slog.String("error", err.Error()))
			return
		}

		h.logger.Info("Order status updated to cancelled",
			slog.Uint64("order_id", uint64(orderCancelledEvent.OrderID)),
			slog.String("reason", orderCancelledEvent.Reason))
	}
}

// PaymentStatusUpdatedEventHandler handles the PaymentStatusUpdatedEvent
type PaymentStatusUpdatedEventHandler struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	ticketRepo  repository.TicketRepository
	eventRepo   repository.EventRepository
	logger      *slog.Logger
}

// NewPaymentStatusUpdatedEventHandler creates a new PaymentStatusUpdatedEventHandler
func NewPaymentStatusUpdatedEventHandler(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	ticketRepo repository.TicketRepository,
	eventRepo repository.EventRepository,
	logger *slog.Logger,
) *PaymentStatusUpdatedEventHandler {
	return &PaymentStatusUpdatedEventHandler{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		ticketRepo:  ticketRepo,
		eventRepo:   eventRepo,
		logger:      logger,
	}
}

// Handle processes the PaymentStatusUpdatedEvent
func (h *PaymentStatusUpdatedEventHandler) Handle(event Event) {
	if paymentStatusEvent, ok := event.(PaymentStatusUpdatedEvent); ok {
		h.logger.Info("Processing payment status update",
			slog.Uint64("payment_id", uint64(paymentStatusEvent.PaymentID)),
			slog.String("status", string(paymentStatusEvent.Status)))

		switch paymentStatusEvent.Status {
		case model.PaymentStatusSuccess:
			// Get order with line items and user
			order, err := h.orderRepo.GetOrderByIDWithLineItems(paymentStatusEvent.OrderID)
			if err != nil {
				h.logger.Error("Failed to get order by ID with line items for status update",
					slog.Uint64("order_id", uint64(paymentStatusEvent.OrderID)),
					slog.String("error", err.Error()))
				return
			}

			// Publish OrderPaidEvent
			orderPaidEvent := OrderPaidEvent{
				OrderID:    order.ID,
				UserID:     order.UserID,
				TotalPrice: order.TotalPrice, // Already in smallest currency unit
				PaidAt:     paymentStatusEvent.UpdatedAt,
			}

			// In a real implementation, you would publish this event to the event bus
			// For now, we'll handle it directly
			orderPaidHandler := NewOrderPaidEventHandler(h.orderRepo, h.logger)
			orderPaidHandler.Handle(orderPaidEvent)

			// Generate tickets
			h.generateTicketsForOrder(order)
		case model.PaymentStatusFailed:
			// Get order
			order, err := h.orderRepo.GetOrderByID(paymentStatusEvent.OrderID)
			if err != nil {
				h.logger.Error("Failed to get order for failed payment",
					slog.Uint64("order_id", uint64(paymentStatusEvent.OrderID)),
					slog.String("error", err.Error()))
				return
			}

			// Restore quotas
			h.restoreQuotasForOrder(order)

			// Publish OrderCancelledEvent
			orderCancelledEvent := OrderCancelledEvent{
				OrderID:     order.ID,
				UserID:      order.UserID,
				Reason:      "Payment failed",
				CancelledAt: paymentStatusEvent.UpdatedAt,
			}

			// In a real implementation, you would publish this event to the event bus
			// For now, we'll handle it directly
			orderCancelledHandler := NewOrderCancelledEventHandler(h.orderRepo, h.logger)
			orderCancelledHandler.Handle(orderCancelledEvent)
		}
	}
}

// generateTicketsForOrder generates tickets for a successful order
func (h *PaymentStatusUpdatedEventHandler) generateTicketsForOrder(order *model.Order) {
	var ticketsToCreate []model.Ticket

	for _, lineItem := range order.OrderLineItems {
		// Fetch EventPrice details to get Type (Name)
		eventPrice, err := h.eventRepo.GetEventPriceByID(lineItem.EventPriceID)
		if err != nil {
			h.logger.Error("failed to get event price for ticket generation",
				slog.Uint64("event_price_id", uint64(lineItem.EventPriceID)),
				slog.String("error", err.Error()))
			continue
		}

		for i := 0; i < lineItem.Quantity; i++ {
			ticket := model.Ticket{
				OrderID:      order.ID,
				EventPriceID: lineItem.EventPriceID,
				Price:        lineItem.PricePerUnit,
				Type:         eventPrice.Name,   // Use EventPrice Name as Ticket Type
				TicketCode:   random.String(10), // Generate unique ticket code
				OwnerName:    order.User.Name,
				OwnerEmail:   order.User.Email,
			}
			ticketsToCreate = append(ticketsToCreate, ticket)
		}
	}

	if len(ticketsToCreate) > 0 {
		if err := h.ticketRepo.CreateTickets(ticketsToCreate); err != nil {
			h.logger.Error("failed to create tickets", slog.String("error", err.Error()))
			return
		}

		// Extract ticket codes for the event
		ticketCodes := make([]string, len(ticketsToCreate))
		for i, ticket := range ticketsToCreate {
			ticketCodes[i] = ticket.TicketCode
		}

		// // Publish TicketsGeneratedEvent
		// ticketsGeneratedEvent := TicketsGeneratedEvent{
		// 	OrderID:     order.ID,
		// 	TicketCodes: ticketCodes,
		// 	GeneratedAt: time.Now(), // Using current time since we don't have the exact timestamp here
		// }

		h.logger.Info("Tickets generated for order",
			slog.Uint64("order_id", uint64(order.ID)),
			slog.Int("count", len(ticketsToCreate)))
	}
}

// restoreQuotasForOrder restores the quotas for a failed order
func (h *PaymentStatusUpdatedEventHandler) restoreQuotasForOrder(order *model.Order) {
	// Get order line items to determine which quotas to restore
	orderWithLineItems, err := h.orderRepo.GetOrderByIDWithLineItems(order.ID)
	if err != nil {
		h.logger.Error("Failed to get order with line items for quota restoration",
			slog.Uint64("order_id", uint64(order.ID)),
			slog.String("error", err.Error()))
		return
	}

	// Restore quotas for each line item
	for _, lineItem := range orderWithLineItems.OrderLineItems {
		err := h.orderRepo.RestoreQuota(lineItem.EventPriceID, lineItem.Quantity)
		if err != nil {
			h.logger.Error("Failed to restore quota",
				slog.Uint64("event_price_id", uint64(lineItem.EventPriceID)),
				slog.Int("quantity", lineItem.Quantity),
				slog.String("error", err.Error()))
		}
	}

	h.logger.Info("Quotas restored for cancelled order",
		slog.Uint64("order_id", uint64(order.ID)))
}
