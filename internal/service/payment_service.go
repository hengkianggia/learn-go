package service

import (
	"errors"
	"fmt"
	"learn/internal/config"
	"learn/internal/dto"
	apperrors "learn/internal/errors"
	midtransGateway "learn/internal/gateway/midtrans"
	"learn/internal/model"
	"learn/internal/pkg/events"
	"learn/internal/pkg/queue"
	"learn/internal/repository"
	"log/slog"
	"strconv"
	"time"

	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(req *dto.CreatePaymentRequest, userID uint) (*model.Payment, error)
	GetPaymentByID(paymentID uint) (*model.Payment, error)
	GetPaymentByOrderID(orderID uint) (*model.Payment, error)
	UpdatePayment(paymentID uint, req *dto.UpdatePaymentRequest) (*model.Payment, error)
	UpdatePaymentStatus(paymentID uint, status model.PaymentStatus) (*model.Payment, error)
	DeletePayment(paymentID uint) error
	HandleNotification(payload map[string]interface{}) error
}

type paymentService struct {
	paymentRepository repository.PaymentRepository
	orderRepository   repository.OrderRepository
	ticketRepository  repository.TicketRepository
	eventRepository   repository.EventRepository
	logger            *slog.Logger
	jobQueue          *queue.JobQueue
	eventBus          *events.EventBus
	midtransGateway   midtransGateway.MidtransGateway // Added
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, ticketRepo repository.TicketRepository, eventRepo repository.EventRepository, logger *slog.Logger, eventBus *events.EventBus) PaymentService {
	return &paymentService{
		paymentRepository: paymentRepo,
		orderRepository:   orderRepo,
		ticketRepository:  ticketRepo,
		eventRepository:   eventRepo,
		logger:            logger,
		jobQueue:          queue.NewJobQueue(5, logger),
		eventBus:          eventBus,
		midtransGateway:   midtransGateway.NewMidtransGateway(logger), // Initialize Gateway
	}
}

func (s *paymentService) CreatePayment(req *dto.CreatePaymentRequest, userID uint) (*model.Payment, error) {
	// Check if order exists
	order, err := s.orderRepository.GetOrderByID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewBusinessRuleError("order_exists", "order not found")
		}
		s.logger.Error("failed to get order by ID", slog.Uint64("order_id", uint64(req.OrderID)), slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("get_order_by_id", err)
	}

	if order.UserID != userID {
		return nil, apperrors.NewBusinessRuleError("payment_authorization", "you are not authorized to pay for this order")
	}

	// Check if a payment already exists for this order
	existingPayment, err := s.paymentRepository.GetPaymentByOrderID(req.OrderID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("failed to check existing payment for order", slog.Uint64("order_id", uint64(req.OrderID)), slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("check_existing_payment", err)
	}

	if existingPayment != nil {
		return nil, apperrors.NewBusinessRuleError("payment_unique", "payment already exists for this order")
	}

	// Process Payment with Midtrans
	var midtransResp *coreapi.ChargeResponse
	var midtransErr error

	// Create unique Order ID for Midtrans (orderID-timestamp)
	midtransOrderID := fmt.Sprintf("ORDER-%d-%d", order.ID, time.Now().Unix())

	switch req.PaymentMethod {
	case model.PaymentMethodBankTransferBCA:
		midtransResp, midtransErr = s.midtransGateway.ChargeBankTransfer(midtransOrderID, order.TotalPrice, "bca")
	case model.PaymentMethodBankTransferBNI:
		midtransResp, midtransErr = s.midtransGateway.ChargeBankTransfer(midtransOrderID, order.TotalPrice, "bni")
	case model.PaymentMethodBankTransferBRI:
		midtransResp, midtransErr = s.midtransGateway.ChargeBankTransfer(midtransOrderID, order.TotalPrice, "bri")
	case model.PaymentMethodGopay:
		midtransResp, midtransErr = s.midtransGateway.ChargeGopay(midtransOrderID, order.TotalPrice)
	case model.PaymentMethodIndomaret:
		midtransResp, midtransErr = s.midtransGateway.ChargeIndomaret(midtransOrderID, order.TotalPrice, "Payment for Order "+strconv.Itoa(int(order.ID)))
	default:
		return nil, apperrors.NewBusinessRuleError("payment_method", "unsupported payment method")
	}

	if midtransErr != nil {
		return nil, apperrors.NewSystemError("midtrans_charge", midtransErr)
	}

	// Parse Midtrans Response
	payment := &model.Payment{
		OrderID:       req.OrderID,
		PaymentMethod: req.PaymentMethod,
		TransactionID: midtransResp.TransactionID,
		PaymentStatus: model.PaymentStatusPending,
	}

	// Extract specific fields based on payment method
	if len(midtransResp.VaNumbers) > 0 {
		payment.VirtualAccountNumber = midtransResp.VaNumbers[0].VANumber
	}
	if len(midtransResp.Actions) > 0 {
		for _, action := range midtransResp.Actions {
			if action.Name == "generate-qr-code" {
				payment.PaymentURL = action.URL
			}
		}
	}
	if midtransResp.PaymentType == "indomaret" {
		payment.PaymentCode = midtransResp.PaymentCode
	}

	if err := s.paymentRepository.CreatePayment(payment); err != nil {
		s.logger.Error("failed to create payment in repository", slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("create_payment", err)
	}

	// Publish PaymentCreatedEvent
	paymentCreatedEvent := events.PaymentCreatedEvent{
		PaymentID: payment.ID,
		OrderID:   req.OrderID,
		Method:    req.PaymentMethod,
		Amount:    order.TotalPrice,
		CreatedAt: time.Now(),
	}
	s.eventBus.Publish(paymentCreatedEvent)

	return payment, nil
}

func (s *paymentService) GetPaymentByID(paymentID uint) (*model.Payment, error) {
	payment, err := s.paymentRepository.GetPaymentByID(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewBusinessRuleError("payment_exists", "payment not found")
		}
		s.logger.Error("failed to get payment by ID", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("get_payment_by_id", err)
	}
	return payment, nil
}

func (s *paymentService) GetPaymentByOrderID(orderID uint) (*model.Payment, error) {
	payment, err := s.paymentRepository.GetPaymentByOrderID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewBusinessRuleError("payment_exists", "payment not found for this order")
		}
		s.logger.Error("failed to get payment by order ID", slog.Uint64("order_id", uint64(orderID)), slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("get_payment_by_order_id", err)
	}
	return payment, nil
}

func (s *paymentService) UpdatePayment(paymentID uint, req *dto.UpdatePaymentRequest) (*model.Payment, error) {
	payment, err := s.paymentRepository.GetPaymentByID(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewBusinessRuleError("payment_exists", "payment not found")
		}
		s.logger.Error("failed to get payment by ID for update", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("get_payment_for_update", err)
	}

	if req.PaymentMethod != nil {
		payment.PaymentMethod = *req.PaymentMethod
	}
	if req.TransactionID != nil {
		payment.TransactionID = *req.TransactionID
	}

	if err := s.paymentRepository.UpdatePayment(payment); err != nil {
		s.logger.Error("failed to update payment in repository", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("update_payment", err)
	}
	return payment, nil
}

func (s *paymentService) UpdatePaymentStatus(paymentID uint, status model.PaymentStatus) (*model.Payment, error) {
	// Use Redis lock to prevent concurrent updates to the same payment
	lockKey := "payment_lock:" + fmt.Sprintf("%d", paymentID)
	set, err := s.paymentRepository.GetRedisClient().SetNX(config.Ctx, lockKey, "locked", 30*time.Second).Result()
	if err != nil {
		s.logger.Error("failed to acquire payment lock", slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("redis_lock", err)
	} else if !set {
		return nil, apperrors.NewBusinessRuleError("payment_processing", "payment is being processed, please wait")
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
			return nil, apperrors.NewBusinessRuleError("payment_exists", "payment not found")
		}

		s.logger.Error("failed to get payment by ID for status update", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("get_payment_for_status_update", err)
	}

	payment.PaymentStatus = status
	if err := s.paymentRepository.UpdatePayment(payment); err != nil {
		s.logger.Error("failed to update payment status in repository", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("update_payment_status", err)
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
			return apperrors.NewBusinessRuleError("payment_exists", "payment not found")
		}
		s.logger.Error("failed to delete payment from repository", slog.Uint64("payment_id", uint64(paymentID)), slog.String("error", err.Error()))
		return apperrors.NewSystemError("delete_payment", err)
	}
	return nil
}
