package service

import (
	"learn/internal/config"
	"learn/internal/dto"
	apperrors "learn/internal/errors"
	"learn/internal/model"
	"learn/internal/pkg/events"
	"learn/internal/repository"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type orderService struct {
	orderRepo repository.OrderRepository
	logger    *slog.Logger
	redis     *redis.Client
	eventBus  *events.EventBus
}

type OrderService interface {
	CreateOrder(input dto.NewOrderInput, userID uint) (*model.Order, error)
}

func NewOrderService(orderRepo repository.OrderRepository, logger *slog.Logger, eventBus *events.EventBus) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		logger:    logger,
		redis:     config.Rdb, // Using the global Redis client from config
		eventBus:  eventBus,
	}
}

func (s *orderService) CreateOrder(input dto.NewOrderInput, userID uint) (*model.Order, error) {
	eventID, err := strconv.ParseUint(input.EventID, 10, 32)
	if err != nil {
		return nil, apperrors.NewValidationError("event_id", "invalid event id", input.EventID)
	}

	event, err := s.orderRepo.GetEventByID(uint(eventID))
	if err != nil {
		return nil, apperrors.NewBusinessRuleError("event_exists", "event not found")
	}

	if event.Status != model.Published {
		return nil, apperrors.NewBusinessRuleError("event_published", "event is not published")
	}

	now := time.Now()
	if now.Before(event.SalesStartDate) || now.After(event.SalesEndDate) {
		return nil, apperrors.NewBusinessRuleError("event_sales_period", "event is not within sales period")
	}

	var priceIDs []uint
	quantityMap := make(map[uint]int)
	totalQuantity := 0 // Track total tickets ordered

	for _, ticketOrder := range input.TicketsOrdered {
		priceID, err := strconv.ParseUint(ticketOrder.PriceId, 10, 32)
		if err != nil {
			return nil, apperrors.NewValidationError("price_id", "invalid price id", ticketOrder.PriceId)
		}

		priceIDs = append(priceIDs, uint(priceID))
		quantityMap[uint(priceID)] = ticketOrder.Quantity
		totalQuantity += ticketOrder.Quantity
	}

	// Anti-abuse validation: Limit total tickets per order
	if totalQuantity > 4 { // Maximum 4 tickets per order
		return nil, apperrors.NewValidationError("tickets_ordered", "maximum 4 tickets allowed per order", totalQuantity)
	}

	// Anti-spam validation: Check if user has too many pending orders recently
	recentOrderCount, err := s.checkRecentOrders(userID)
	if err != nil {
		s.logger.Error("failed to check recent orders", slog.String("error", err.Error()))
	}
	if recentOrderCount > 1 { // Allow max 5 orders per hour
		return nil, apperrors.NewBusinessRuleError("order_limit", "too many orders recently, please wait before placing another order")
	}

	prices, err := s.orderRepo.GetEventPricesByIDs(priceIDs)
	if err != nil {
		return nil, apperrors.NewBusinessRuleError("event_prices_exist", "failed to get prices")
	}

	if len(prices) != len(priceIDs) {
		return nil, apperrors.NewBusinessRuleError("event_prices_exist", "one or more prices not found")
	}

	var totalPrice int64
	priceUpdates := make(map[uint]int)

	for _, price := range prices {
		if price.EventID != uint(eventID) {
			return nil, apperrors.NewBusinessRuleError("event_prices_match", "one or more prices do not belong to this event")
		}

		quantity := quantityMap[price.ID]
		if price.Quota < quantity {
			return nil, apperrors.NewBusinessRuleError("ticket_quota", "not enough quota for ticket")
		}

		// Calculate total price using integer arithmetic to avoid floating point errors
		itemTotal := price.Price * int64(quantity)

		// Check for overflow before adding
		const maxInt64 = int64(^uint64(0) >> 1)
		if totalPrice > (maxInt64 - itemTotal) {
			return nil, apperrors.NewBusinessRuleError("order_total_limit", "order total exceeds maximum allowed value")
		}

		totalPrice += itemTotal
		priceUpdates[price.ID] = quantity
	}

	// Anti-double spending validation: Check if user is trying to order same tickets again
	orderLockKey := "order_lock:" + strconv.FormatUint(uint64(userID), 10) + ":" + input.EventID
	set, err := s.redis.SetNX(config.Ctx, orderLockKey, "locked", 5*time.Minute).Result()

	if err != nil {
		s.logger.Error("failed to set order lock", slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("redis_lock", err)
	} else if !set {
		return nil, apperrors.NewBusinessRuleError("order_lock", "order is being processed, please wait")
	}

	defer s.redis.Del(config.Ctx, orderLockKey) // Clean up lock

	order := &model.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     model.OrderPending,
		PaymentDue: time.Now().Add(24 * time.Hour),
	}

	err = s.orderRepo.CreateOrderInTransaction(order, prices, priceUpdates)
	if err != nil {
		s.logger.Error("failed to create order", slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("create_order_transaction", err)
	}

	// Publish OrderCreatedEvent
	orderCreatedEvent := events.OrderCreatedEvent{
		OrderID:    order.ID,
		UserID:     userID,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
	}
	s.eventBus.Publish(orderCreatedEvent)

	return order, nil
}

// checkRecentOrders checks how many orders a user has placed in the last hour
func (s *orderService) checkRecentOrders(userID uint) (int, error) {
	key := "recent_orders:" + strconv.FormatUint(uint64(userID), 10)
	count, err := s.redis.Incr(config.Ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration for the counter to 1 hour
	if count == 1 { // First increment, set TTL
		s.redis.Expire(config.Ctx, key, time.Hour)
	}

	return int(count), nil
}
