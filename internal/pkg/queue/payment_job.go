package queue

import (
	"fmt"
	"learn/internal/model"
	"learn/internal/pkg/random"
	"learn/internal/repository"
	"log/slog"
)

// PaymentJob represents a job for processing payment updates
type PaymentJob struct {
	PaymentID uint
	Status    model.PaymentStatus
	OrderRepo repository.OrderRepository
	PaymentRepo repository.PaymentRepository
	TicketRepo repository.TicketRepository
	EventRepo repository.EventRepository
	Logger *slog.Logger
}

// NewPaymentJob creates a new payment job
func NewPaymentJob(
	paymentID uint,
	status model.PaymentStatus,
	orderRepo repository.OrderRepository,
	paymentRepo repository.PaymentRepository,
	ticketRepo repository.TicketRepository,
	eventRepo repository.EventRepository,
	logger *slog.Logger,
) Job {
	paymentJob := &PaymentJob{
		PaymentID: paymentID,
		Status: status,
		OrderRepo: orderRepo,
		PaymentRepo: paymentRepo,
		TicketRepo: ticketRepo,
		EventRepo: eventRepo,
		Logger: logger,
	}

	return Job{
		ID:       "payment_update_" + fmt.Sprintf("%d", paymentID),
		Task:     paymentJob.execute,
		RetryMax: 3,
		Retry:    0,
	}
}

// execute contains the actual business logic for processing payment
func (pj *PaymentJob) execute() error {
	payment, err := pj.PaymentRepo.GetPaymentByID(pj.PaymentID)
	if err != nil {
		pj.Logger.Error("Failed to get payment by ID",
			slog.Uint64("payment_id", uint64(pj.PaymentID)),
			slog.String("error", err.Error()))
		return err
	}

	// Update payment status
	payment.PaymentStatus = pj.Status
	if err := pj.PaymentRepo.UpdatePayment(payment); err != nil {
		pj.Logger.Error("Failed to update payment status in repository",
			slog.Uint64("payment_id", uint64(pj.PaymentID)),
			slog.String("error", err.Error()))
		return err
	}

	if pj.Status == model.PaymentStatusSuccess {
		// Get order with line items and user
		order, err := pj.OrderRepo.GetOrderByIDWithLineItems(payment.OrderID)
		if err != nil {
			pj.Logger.Error("Failed to get order by ID with line items for status update",
				slog.Uint64("order_id", uint64(payment.OrderID)),
				slog.String("error", err.Error()))
			return err
		}

		// Update order status to paid
		order.Status = model.OrderPaid
		if err := pj.OrderRepo.UpdateOrder(order); err != nil {
			pj.Logger.Error("Failed to update order status in repository",
				slog.Uint64("order_id", uint64(payment.OrderID)),
				slog.String("error", err.Error()))
			return err
		}

		// Generate and create tickets directly without service dependency
		err = pj.generateTicketsForOrder(order)
		if err != nil {
			pj.Logger.Error("Failed to generate tickets for order",
				slog.Uint64("order_id", uint64(order.ID)),
				slog.String("error", err.Error()))
			return err
		}
	} else if pj.Status == model.PaymentStatusFailed {
		// Handle failed payment - restore quotas
		order, err := pj.OrderRepo.GetOrderByID(payment.OrderID)
		if err != nil {
			pj.Logger.Error("Failed to get order for failed payment",
				slog.Uint64("order_id", uint64(payment.OrderID)),
				slog.String("error", err.Error()))
			return err
		}

		// Restore quotas for the order
		err = pj.restoreQuotasForOrder(order)
		if err != nil {
			pj.Logger.Error("Failed to restore quotas for failed order",
				slog.Uint64("order_id", uint64(order.ID)),
				slog.String("error", err.Error()))
			return err
		}

		// Update order status to cancelled
		order.Status = model.OrderCancelled
		if err := pj.OrderRepo.UpdateOrder(order); err != nil {
			pj.Logger.Error("Failed to update order status to cancelled",
				slog.Uint64("order_id", uint64(order.ID)),
				slog.String("error", err.Error()))
			return err
		}
	}

	return nil
}

// generateTicketsForOrder generates tickets for a successful order
func (pj *PaymentJob) generateTicketsForOrder(order *model.Order) error {
	var ticketsToCreate []model.Ticket

	for _, lineItem := range order.OrderLineItems {
		// Fetch EventPrice details to get Type (Name)
		eventPrice, err := pj.EventRepo.GetEventPriceByID(lineItem.EventPriceID)
		if err != nil {
			pj.Logger.Error("failed to get event price for ticket generation",
				slog.Uint64("event_price_id", uint64(lineItem.EventPriceID)),
				slog.String("error", err.Error()))
			return err
		}

		for i := 0; i < lineItem.Quantity; i++ {
			ticketsToCreate = append(ticketsToCreate, model.Ticket{
				OrderID:      order.ID,
				EventPriceID: lineItem.EventPriceID,
				Price:        lineItem.PricePerUnit,
				Type:         eventPrice.Name, // Use EventPrice Name as Ticket Type
				TicketCode:   random.String(10),
				OwnerName:    order.User.Name,
				OwnerEmail:   order.User.Email,
			})
		}
	}

	if len(ticketsToCreate) > 0 {
		if err := pj.TicketRepo.CreateTickets(ticketsToCreate); err != nil {
			pj.Logger.Error("failed to create tickets", slog.String("error", err.Error()))
			return err
		}
	}

	return nil
}

// restoreQuotasForOrder restores the quotas for a failed order
func (pj *PaymentJob) restoreQuotasForOrder(order *model.Order) error {
	// Get order line items to determine which quotas to restore
	orderWithLineItems, err := pj.OrderRepo.GetOrderByIDWithLineItems(order.ID)
	if err != nil {
		return err
	}

	// Restore quotas for each line item
	for _, lineItem := range orderWithLineItems.OrderLineItems {
		err := pj.OrderRepo.RestoreQuota(lineItem.EventPriceID, lineItem.Quantity)
		if err != nil {
			pj.Logger.Error("Failed to restore quota",
				slog.Uint64("event_price_id", uint64(lineItem.EventPriceID)),
				slog.Int("quantity", lineItem.Quantity),
				slog.String("error", err.Error()))
			return err
		}
	}

	return nil
}