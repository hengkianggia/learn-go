package scheduler

import (
	"learn/internal/model"
	"learn/internal/repository"
	"learn/internal/service"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// OrderExpirationScheduler handles automatic cancellation of expired orders
type OrderExpirationScheduler struct {
	orderRepo repository.OrderRepository
	logger    *slog.Logger
	stopChan  chan struct{}
}

// NewOrderExpirationScheduler creates a new order expiration scheduler
func NewOrderExpirationScheduler(orderRepo repository.OrderRepository, logger *slog.Logger) *OrderExpirationScheduler {
	return &OrderExpirationScheduler{
		orderRepo: orderRepo,
		logger:    logger,
		stopChan:  make(chan struct{}),
	}
}

// Start begins the scheduled task to check for expired orders
func (s *OrderExpirationScheduler) Start() {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	s.logger.Info("Order expiration scheduler started")

	for {
		select {
		case <-ticker.C:
			s.checkExpiredOrders()
		case <-s.stopChan:
			s.logger.Info("Order expiration scheduler stopped")
			return
		}
	}
}

// Stop stops the scheduled task
func (s *OrderExpirationScheduler) Stop() {
	close(s.stopChan)
}

// checkExpiredOrders finds and cancels orders that have exceeded their payment due time
func (s *OrderExpirationScheduler) checkExpiredOrders() {
	now := time.Now()

	// Find orders that are pending and have exceeded their payment due time
	var expiredOrders []model.Order
	err := s.orderRepo.GetDB().Where("status = ? AND payment_due < ?", model.OrderPending, now).Find(&expiredOrders).Error
	if err != nil {
		s.logger.Error("Failed to find expired orders", slog.String("error", err.Error()))
		return
	}

	if len(expiredOrders) == 0 {
		s.logger.Debug("No expired orders found")
		return
	}

	s.logger.Info("Found expired orders", slog.Int("count", len(expiredOrders)))

	for _, order := range expiredOrders {
		err := s.cancelOrder(&order, "Expired - Payment not received within deadline")
		if err != nil {
			s.logger.Error("Failed to cancel expired order",
				slog.Uint64("order_id", uint64(order.ID)),
				slog.String("error", err.Error()))
		} else {
			s.logger.Info("Successfully cancelled expired order",
				slog.Uint64("order_id", uint64(order.ID)))
		}
	}
}

// cancelOrder cancels an order and restores the quotas
func (s *OrderExpirationScheduler) cancelOrder(order *model.Order, reason string) error {
	// Update order status to cancelled
	order.Status = model.OrderCancelled
	err := s.orderRepo.UpdateOrder(order)
	if err != nil {
		return err
	}

	// Restore quotas for the order
	err = s.restoreQuotasForOrder(order)
	if err != nil {
		s.logger.Error("Failed to restore quotas for cancelled order",
			slog.Uint64("order_id", uint64(order.ID)),
			slog.String("error", err.Error()))
		// Don't return error here as the order is already cancelled
	}

	return nil
}

// restoreQuotasForOrder restores the quotas for a cancelled order
func (s *OrderExpirationScheduler) restoreQuotasForOrder(order *model.Order) error {
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