package service

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
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
	logger            *slog.Logger
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, logger *slog.Logger) PaymentService {
	return &paymentService{
		paymentRepository: paymentRepo,
		orderRepository:   orderRepo,
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
		order, err := s.orderRepository.GetOrderByID(payment.OrderID)
		if err != nil {
			s.logger.Error("failed to get order by ID for status update", slog.Uint64("order_id", uint64(payment.OrderID)), slog.String("error", err.Error()))
			return nil, errors.New("failed to update order status")
		}
		order.Status = model.OrderPaid
		if err := s.orderRepository.UpdateOrder(order); err != nil {
			s.logger.Error("failed to update order status in repository", slog.Uint64("order_id", uint64(payment.OrderID)), slog.String("error", err.Error()))
			return nil, errors.New("failed to update order status")
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
