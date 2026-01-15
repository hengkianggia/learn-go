package service

import (
	apperrors "learn/internal/errors"
	"learn/internal/model"
	"learn/internal/pkg/events"
	"learn/internal/repository"
	"log/slog"
)

type orderCancellationService struct {
	orderRepo repository.OrderRepository
	logger    *slog.Logger
	eventBus  *events.EventBus
}

type OrderCancellationService interface {
	CancelOrder(orderID uint, userID uint, reason string) error
}

func NewOrderCancellationService(orderRepo repository.OrderRepository, logger *slog.Logger, eventBus *events.EventBus) OrderCancellationService {
	return &orderCancellationService{
		orderRepo: orderRepo,
		logger:    logger,
		eventBus:  eventBus,
	}
}

func (s *orderCancellationService) CancelOrder(orderID uint, userID uint, reason string) error {
	// Get the order
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		s.logger.Error("Failed to get order for cancellation",
			slog.Uint64("order_id", uint64(orderID)),
			slog.String("error", err.Error()))
		return apperrors.NewBusinessRuleError("order_exists", "order not found")
	}

	// Check if the order belongs to the user
	if order.UserID != userID {
		s.logger.Error("User not authorized to cancel order",
			slog.Uint64("user_id", uint64(userID)),
			slog.Uint64("order_id", uint64(orderID)))
		return apperrors.NewBusinessRuleError("order_authorization", "you are not authorized to cancel this order")
	}

	// Check if order is already cancelled or paid
	if order.Status == model.OrderCancelled {
		return apperrors.NewBusinessRuleError("order_status", "order is already cancelled")
	}
	if order.Status == model.OrderPaid {
		return apperrors.NewBusinessRuleError("order_status", "paid orders cannot be cancelled manually")
	}

	// Update order status to cancelled
	order.Status = model.OrderCancelled
	err = s.orderRepo.UpdateOrder(order)
	if err != nil {
		s.logger.Error("Failed to update order status to cancelled",
			slog.Uint64("order_id", uint64(orderID)),
			slog.String("error", err.Error()))
		return apperrors.NewSystemError("update_order_status", err)
	}

	// Restore quotas for the order
	err = s.restoreQuotasForOrder(order)
	if err != nil {
		s.logger.Error("Failed to restore quotas for cancelled order",
			slog.Uint64("order_id", uint64(orderID)),
			slog.String("error", err.Error()))
		// Don't return error here as the order is already cancelled
	}

	// Publish OrderCancelledEvent
	orderCancelledEvent := events.OrderCancelledEvent{
		OrderID:     order.ID,
		UserID:      userID,
		Reason:      reason,
		CancelledAt: order.UpdatedAt,
	}
	s.eventBus.Publish(orderCancelledEvent)

	s.logger.Info("Order cancelled successfully",
		slog.Uint64("order_id", uint64(orderID)),
		slog.String("reason", reason))

	return nil
}

// restoreQuotasForOrder restores the quotas for a cancelled order
func (s *orderCancellationService) restoreQuotasForOrder(order *model.Order) error {
	// Get order line items to determine which quotas to restore
	orderWithLineItems, err := s.orderRepo.GetOrderByIDWithLineItems(order.ID)
	if err != nil {
		return err
	}

	// Restore quotas for each line item
	for _, lineItem := range orderWithLineItems.OrderLineItems {
		err := s.orderRepo.RestoreQuota(lineItem.EventPriceID, lineItem.Quantity)
		if err != nil {
			s.logger.Error("Failed to restore quota",
				slog.Uint64("event_price_id", uint64(lineItem.EventPriceID)),
				slog.Int("quantity", lineItem.Quantity),
				slog.String("error", err.Error()))
			return err
		}
	}

	s.logger.Info("Quotas restored for cancelled order",
		slog.Uint64("order_id", uint64(order.ID)))

	return nil
}
