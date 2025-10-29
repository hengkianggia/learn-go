package service

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/random" // Added
	"learn/internal/repository"
	"log/slog"

	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(req *dto.CreatePaymentRequest, userID uint) (*model.Payment, error)
	GetPaymentByID(paymentID uint) (*model.Payment, error)
	GetPaymentByOrderID(orderID uint) (*model.Payment, error)
	UpdatePayment(paymentID uint, req *dto.UpdatePaymentRequest) (*model.Payment, error)
	UpdatePaymentStatus(paymentID uint, status model.PaymentStatus) (*model.Payment, error)
	DeletePayment(paymentID uint) error
}

type paymentService struct {
	paymentRepository repository.PaymentRepository
	orderRepository   repository.OrderRepository
	ticketRepository  repository.TicketRepository // Added
	eventRepository   repository.EventRepository  // Added
	logger            *slog.Logger
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, ticketRepo repository.TicketRepository, eventRepo repository.EventRepository, logger *slog.Logger) PaymentService {
	return &paymentService{
		paymentRepository: paymentRepo,
		orderRepository:   orderRepo,
		ticketRepository:  ticketRepo,
		eventRepository:   eventRepo,
		logger:            logger,
	}
}

func (s *paymentService) CreatePayment(req *dto.CreatePaymentRequest, userID uint) (*model.Payment, error) {
	// Check if order exists
	order, err := s.orderRepository.GetOrderByID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		s.logger.Error("failed to get order by ID", slog.Uint64("order_id", uint64(req.OrderID)), slog.String("error", err.Error()))
		return nil, errors.New("failed to create payment")
	}

	if order.UserID != userID {
		return nil, errors.New("you are not authorized to pay for this order")
	}

	// Check if a payment already exists for this order (due to 1:1 relationship)
	existingPayment, err := s.paymentRepository.GetPaymentByOrderID(req.OrderID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("failed to check existing payment for order", slog.Uint64("order_id", uint64(req.OrderID)), slog.String("error", err.Error()))
		return nil, errors.New("failed to create payment")
	}

	if existingPayment != nil {
		return nil, errors.New("payment already exists for this order")
	}

	// Validate amount against order total price (business logic)
	// Assuming TotalPrice is float64 and converting to smallest currency unit (e.g., cents)
	// This conversion should be handled carefully to avoid precision issues.
	// For example, if TotalPrice is 100.50, it should be 10050.
	// A safer approach might be to store TotalPrice as int64 in Order model as well.
	payment := &model.Payment{
		OrderID:       req.OrderID,
		PaymentMethod: req.PaymentMethod,
		TransactionID: req.TransactionID,
		// Amount:        order.TotalPrice,
		PaymentStatus: model.PaymentStatusPending,
	}

	if err := s.paymentRepository.CreatePayment(payment); err != nil {
		s.logger.Error("failed to create payment in repository", slog.String("error", err.Error()))
		return nil, errors.New("failed to create payment")
	}
	return payment, nil
}

func (s *paymentService) GetPaymentByID(paymentID uint) (*model.Payment, error) {
	payment, err := s.paymentRepository.GetPaymentByID(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		s.logger.Error("failed to get payment by ID", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, errors.New("failed to retrieve payment")
	}
	return payment, nil
}

func (s *paymentService) GetPaymentByOrderID(orderID uint) (*model.Payment, error) {
	payment, err := s.paymentRepository.GetPaymentByOrderID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found for this order")
		}
		s.logger.Error("failed to get payment by order ID", slog.Uint64("order_id", uint64(orderID)), slog.String("error", err.Error()))
		return nil, errors.New("failed to retrieve payment for order")
	}
	return payment, nil
}

func (s *paymentService) UpdatePayment(paymentID uint, req *dto.UpdatePaymentRequest) (*model.Payment, error) {
	payment, err := s.paymentRepository.GetPaymentByID(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		s.logger.Error("failed to get payment by ID for update", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, errors.New("failed to update payment")
	}

	if req.PaymentMethod != nil {
		payment.PaymentMethod = *req.PaymentMethod
	}
	if req.TransactionID != nil {
		payment.TransactionID = *req.TransactionID
	}

	if err := s.paymentRepository.UpdatePayment(payment); err != nil {
		s.logger.Error("failed to update payment in repository", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, errors.New("failed to update payment")
	}
	return payment, nil
}

func (s *paymentService) UpdatePaymentStatus(paymentID uint, status model.PaymentStatus) (*model.Payment, error) {
	payment, err := s.paymentRepository.GetPaymentByID(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}

		s.logger.Error("failed to get payment by ID for status update", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, errors.New("failed to update payment status")
	}

	payment.PaymentStatus = status

	if err := s.paymentRepository.UpdatePayment(payment); err != nil {
		s.logger.Error("failed to update payment status in repository", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, errors.New("failed to update payment status")
	}

	if status == model.PaymentStatusSuccess {
		// Get order with line items and user
		order, err := s.orderRepository.GetOrderByIDWithLineItems(payment.OrderID)
		if err != nil {
			s.logger.Error("failed to get order by ID with line items for status update", slog.Uint64("order_id", uint64(payment.OrderID)), slog.String("error", err.Error()))
			return nil, errors.New("failed to update order status")
		}

		// Update order status to paid
		order.Status = model.OrderPaid
		if err := s.orderRepository.UpdateOrder(order); err != nil {
			s.logger.Error("failed to update order status in repository", slog.Uint64("order_id", uint64(payment.OrderID)), slog.String("error", err.Error()))
			return nil, errors.New("failed to update order status")
		}

		// Generate and create tickets
		var ticketsToCreate []model.Ticket
		for _, lineItem := range order.OrderLineItems {
			// Fetch EventPrice details to get Type (Name)
			eventPrice, err := s.eventRepository.GetEventPriceByID(lineItem.EventPriceID)
			if err != nil {
				s.logger.Error("failed to get event price for ticket generation", slog.Uint64("event_price_id", uint64(lineItem.EventPriceID)), slog.String("error", err.Error()))
				return nil, errors.New("failed to generate tickets")
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
			if err := s.ticketRepository.CreateTickets(ticketsToCreate); err != nil {
				s.logger.Error("failed to create tickets", slog.String("error", err.Error()))
				return nil, errors.New("failed to generate tickets")
			}
		}
	}

	return payment, nil
}

func (s *paymentService) DeletePayment(paymentID uint) error {
	if err := s.paymentRepository.DeletePayment(paymentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payment not found")
		}
		s.logger.Error("failed to delete payment from repository", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return errors.New("failed to delete payment")
	}
	return nil
}
