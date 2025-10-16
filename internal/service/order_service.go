package service

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/repository"
	"log/slog"
	"strconv"
	"time"
)

type OrderService interface {
	CreateOrder(input dto.NewOrderInput, userID uint) (*model.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	logger    *slog.Logger
}

func NewOrderService(orderRepo repository.OrderRepository, logger *slog.Logger) OrderService {
	return &orderService{orderRepo: orderRepo, logger: logger}
}

func (s *orderService) CreateOrder(input dto.NewOrderInput, userID uint) (*model.Order, error) {
	eventID, err := strconv.ParseUint(input.EventID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid event id")
	}

	event, err := s.orderRepo.GetEventByID(uint(eventID))
	if err != nil {
		return nil, errors.New("event not found")
	}

	if event.Status != model.Published {
		return nil, errors.New("event is not published")
	}

	now := time.Now()
	if now.Before(event.SalesStartDate) || now.After(event.SalesEndDate) {
		return nil, errors.New("event is not within sales period")
	}

	var priceIDs []uint
	quantityMap := make(map[uint]int)
	for _, ticketOrder := range input.TicketsOrdered {
		priceID, err := strconv.ParseUint(ticketOrder.PriceId, 10, 32)
		if err != nil {
			return nil, errors.New("invalid price id")
		}
		priceIDs = append(priceIDs, uint(priceID))
		quantityMap[uint(priceID)] = ticketOrder.Quantity
	}

	prices, err := s.orderRepo.GetEventPricesByIDs(priceIDs)
	if err != nil {
		return nil, errors.New("failed to get prices")
	}

	if len(prices) != len(priceIDs) {
		return nil, errors.New("one or more prices not found")
	}

	var totalPrice float64
	var tickets []model.Ticket
	priceUpdates := make(map[uint]int)

	for _, price := range prices {
		if price.EventID != uint(eventID) {
			return nil, errors.New("one or more prices do not belong to this event")
		}

		quantity := quantityMap[price.ID]
		if price.Quota < quantity {
			return nil, errors.New("not enough quota for ticket")
		}

		totalPrice += float64(price.Price * quantity)
		priceUpdates[price.ID] = quantity

		for i := 0; i < quantity; i++ {
			tickets = append(tickets, model.Ticket{
				EventPriceID: price.ID,
				Price:        float64(price.Price),
				Type:         price.Name,
			})
		}
	}

	order := &model.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     model.OrderPending,
		PaymentDue: time.Now().Add(24 * time.Hour),
	}

	err = s.orderRepo.CreateOrderInTransaction(order, tickets, priceUpdates)
	if err != nil {
		s.logger.Error("failed to create order", slog.String("error", err.Error()))
		return nil, err
	}

	return order, nil
}
