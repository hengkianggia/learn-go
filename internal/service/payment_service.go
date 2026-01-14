package service

import (
	"errors"
	"fmt"
	"learn/internal/config"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/events"
	"learn/internal/pkg/queue"
	"learn/internal/repository"
	"log/slog"
	"time"

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
	jobQueue          *queue.JobQueue
	eventBus          *events.EventBus
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, ticketRepo repository.TicketRepository, eventRepo repository.EventRepository, logger *slog.Logger, eventBus *events.EventBus) PaymentService {
	return &paymentService{
		paymentRepository: paymentRepo,
		orderRepository:   orderRepo,
		ticketRepository:  ticketRepo,
		eventRepository:   eventRepo,
		logger:            logger,
		jobQueue:          queue.NewJobQueue(5, logger), // 5 workers for payment processing
		eventBus:          eventBus,
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

	// Create payment record with order total price (already in smallest currency unit)
	payment := &model.Payment{
		OrderID:       req.OrderID,
		PaymentMethod: req.PaymentMethod,
		TransactionID: req.TransactionID,
		PaymentStatus: model.PaymentStatusPending,
	}

	if err := s.paymentRepository.CreatePayment(payment); err != nil {
		s.logger.Error("failed to create payment in repository", slog.String("error", err.Error()))
		return nil, errors.New("failed to create payment")
	}

	// Publish PaymentCreatedEvent
	paymentCreatedEvent := events.PaymentCreatedEvent{
		PaymentID: payment.ID,
		OrderID:   req.OrderID,
		Method:    req.PaymentMethod,
		Amount:    order.TotalPrice, // Using order total as payment amount
		CreatedAt: time.Now(),
	}
	s.eventBus.Publish(paymentCreatedEvent)

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
	// Use Redis lock to prevent concurrent updates to the same payment
	lockKey := "payment_lock:" + fmt.Sprintf("%d", paymentID)
	set, err := s.paymentRepository.GetRedisClient().SetNX(config.Ctx, lockKey, "locked", 30*time.Second).Result()
	if err != nil {
		s.logger.Error("failed to acquire payment lock", slog.String("error", err.Error()))
	} else if !set {
		return nil, errors.New("payment is being processed, please wait")
	}
	defer s.paymentRepository.GetRedisClient().Del(config.Ctx, lockKey) // Clean up lock

	// Create a job for processing the payment status update
	job := queue.NewPaymentJob(
		paymentID,
		status,
		s.orderRepository,
		s.paymentRepository,
		s.ticketRepository,
		s.eventRepository,
		s.logger,
	)

	// Add the job to the queue for asynchronous processing
	s.jobQueue.Enqueue(job)

	// Return the payment with the updated status without waiting for the job to complete
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

	// Publish PaymentStatusUpdatedEvent
	paymentStatusEvent := events.PaymentStatusUpdatedEvent{
		PaymentID: paymentID,
		OrderID:   payment.OrderID,
		Status:    status,
		UpdatedAt: time.Now(),
	}
	s.eventBus.Publish(paymentStatusEvent)

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
