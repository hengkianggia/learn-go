package controller

import (
	"learn/internal/dto"
	apperrors "learn/internal/errors"
	"learn/internal/model"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type paymentController struct {
	paymentService service.PaymentService
	logger         *slog.Logger
}

type PaymentController interface {
	CreatePayment(c *gin.Context)
	GetPaymentByID(c *gin.Context)
	GetPaymentByOrderID(c *gin.Context)
	UpdatePayment(c *gin.Context)
	UpdatePaymentStatus(c *gin.Context)
	DeletePayment(c *gin.Context)
}

func NewPaymentController(paymentService service.PaymentService, logger *slog.Logger) PaymentController {
	return &paymentController{
		paymentService: paymentService,
		logger:         logger,
	}
}

func (ctrl *paymentController) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	user, exists := c.Get("user")
	if !exists {
		response.SendUnauthorizedError(c, "User not authenticated")
		return
	}

	payment, err := ctrl.paymentService.CreatePayment(&req, user.(model.User).ID)
	if err != nil {
		// Handle different types of errors appropriately
		switch appErr := err.(type) {
		case apperrors.ValidationError:
			ctrl.logger.Info("Validation error in create payment",
				slog.String("field", appErr.Field),
				slog.String("message", appErr.Message),
				slog.Any("value", appErr.Value))
			response.SendBadRequestError(c, appErr.Error())
			return
		case apperrors.BusinessRuleError:
			ctrl.logger.Info("Business rule error in create payment",
				slog.String("rule", appErr.Rule),
				slog.String("message", appErr.Message))
			response.SendBadRequestError(c, appErr.Error())
			return
		case apperrors.SystemError:
			ctrl.logger.Error("System error in create payment",
				slog.String("operation", appErr.Operation),
				slog.String("message", appErr.Message),
				slog.Any("error", appErr.Err))
			response.SendInternalServerError(c, ctrl.logger, appErr)
			return
		default:
			// For any other error types, treat as internal server error
			ctrl.logger.Error("Unknown error in create payment", slog.String("error", err.Error()))
			response.SendInternalServerError(c, ctrl.logger, err)
			return
		}
	}

	// Create the response DTO
	paymentResponse := dto.PaymentResponse{
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		PaymentMethod: payment.PaymentMethod,
		TransactionID: payment.TransactionID,
		PaymentStatus: payment.PaymentStatus,
		PaymentDate:   payment.PaymentDate,
	}

	response.SendSuccess(c, http.StatusCreated, "Payment created successfully", paymentResponse)
}

func (ctrl *paymentController) GetPaymentByID(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 32)
	if err != nil {
		response.SendBadRequestError(c, "Invalid payment ID")
		return
	}

	payment, err := ctrl.paymentService.GetPaymentByID(uint(paymentID))
	if err != nil {
		// Handle different types of errors appropriately
		switch appErr := err.(type) {
		case apperrors.ValidationError:
			ctrl.logger.Info("Validation error in get payment by ID",
				slog.String("field", appErr.Field),
				slog.String("message", appErr.Message),
				slog.Any("value", appErr.Value))
			response.SendBadRequestError(c, appErr.Error())
			return
		case apperrors.BusinessRuleError:
			ctrl.logger.Info("Business rule error in get payment by ID",
				slog.String("rule", appErr.Rule),
				slog.String("message", appErr.Message))
			response.SendNotFoundError(c, appErr.Error())
			return
		case apperrors.SystemError:
			ctrl.logger.Error("System error in get payment by ID",
				slog.String("operation", appErr.Operation),
				slog.String("message", appErr.Message),
				slog.Any("error", appErr.Err))
			response.SendInternalServerError(c, ctrl.logger, appErr)
			return
		default:
			// For any other error types, treat as internal server error
			ctrl.logger.Error("Unknown error in get payment by ID", slog.String("error", err.Error()))
			response.SendInternalServerError(c, ctrl.logger, err)
			return
		}
	}

	// Create the response DTO
	paymentResponse := dto.PaymentResponse{
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		PaymentMethod: payment.PaymentMethod,
		TransactionID: payment.TransactionID,
		PaymentStatus: payment.PaymentStatus,
		PaymentDate:   payment.PaymentDate,
	}

	response.SendSuccess(c, http.StatusOK, "Payment retrieved successfully", paymentResponse)
}

func (ctrl *paymentController) GetPaymentByOrderID(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		response.SendBadRequestError(c, "Invalid order ID")
		return
	}

	payment, err := ctrl.paymentService.GetPaymentByOrderID(uint(orderID))
	if err != nil {
		// Handle different types of errors appropriately
		switch appErr := err.(type) {
		case apperrors.ValidationError:
			ctrl.logger.Info("Validation error in get payment by order ID",
				slog.String("field", appErr.Field),
				slog.String("message", appErr.Message),
				slog.Any("value", appErr.Value))
			response.SendBadRequestError(c, appErr.Error())
			return
		case apperrors.BusinessRuleError:
			ctrl.logger.Info("Business rule error in get payment by order ID",
				slog.String("rule", appErr.Rule),
				slog.String("message", appErr.Message))
			response.SendNotFoundError(c, appErr.Error())
			return
		case apperrors.SystemError:
			ctrl.logger.Error("System error in get payment by order ID",
				slog.String("operation", appErr.Operation),
				slog.String("message", appErr.Message),
				slog.Any("error", appErr.Err))
			response.SendInternalServerError(c, ctrl.logger, appErr)
			return
		default:
			// For any other error types, treat as internal server error
			ctrl.logger.Error("Unknown error in get payment by order ID", slog.String("error", err.Error()))
			response.SendInternalServerError(c, ctrl.logger, err)
			return
		}
	}

	// Create the response DTO
	paymentResponse := dto.PaymentResponse{
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		PaymentMethod: payment.PaymentMethod,
		TransactionID: payment.TransactionID,
		PaymentStatus: payment.PaymentStatus,
		PaymentDate:   payment.PaymentDate,
	}

	response.SendSuccess(c, http.StatusOK, "Payment retrieved successfully for order", paymentResponse)
}

func (ctrl *paymentController) UpdatePayment(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 32)
	if err != nil {
		response.SendBadRequestError(c, "Invalid payment ID")
		return
	}

	var req dto.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	payment, err := ctrl.paymentService.UpdatePayment(uint(paymentID), &req)
	if err != nil {
		// Handle different types of errors appropriately
		switch appErr := err.(type) {
		case apperrors.ValidationError:
			ctrl.logger.Info("Validation error in update payment",
				slog.String("field", appErr.Field),
				slog.String("message", appErr.Message),
				slog.Any("value", appErr.Value))
			response.SendBadRequestError(c, appErr.Error())
			return
		case apperrors.BusinessRuleError:
			ctrl.logger.Info("Business rule error in update payment",
				slog.String("rule", appErr.Rule),
				slog.String("message", appErr.Message))
			response.SendNotFoundError(c, appErr.Error())
			return
		case apperrors.SystemError:
			ctrl.logger.Error("System error in update payment",
				slog.String("operation", appErr.Operation),
				slog.String("message", appErr.Message),
				slog.Any("error", appErr.Err))
			response.SendInternalServerError(c, ctrl.logger, appErr)
			return
		default:
			// For any other error types, treat as internal server error
			ctrl.logger.Error("Unknown error in update payment", slog.String("error", err.Error()))
			response.SendInternalServerError(c, ctrl.logger, err)
			return
		}
	}

	response.SendSuccess(c, http.StatusOK, "Payment updated successfully", payment)
}

func (ctrl *paymentController) UpdatePaymentStatus(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 32)
	if err != nil {
		response.SendBadRequestError(c, "Invalid payment ID")
		return
	}

	var req dto.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	payment, err := ctrl.paymentService.UpdatePaymentStatus(uint(paymentID), req.Status)
	if err != nil {
		// Handle different types of errors appropriately
		switch appErr := err.(type) {
		case apperrors.ValidationError:
			ctrl.logger.Info("Validation error in update payment status",
				slog.String("field", appErr.Field),
				slog.String("message", appErr.Message),
				slog.Any("value", appErr.Value))
			response.SendBadRequestError(c, appErr.Error())
			return
		case apperrors.BusinessRuleError:
			ctrl.logger.Info("Business rule error in update payment status",
				slog.String("rule", appErr.Rule),
				slog.String("message", appErr.Message))
			response.SendNotFoundError(c, appErr.Error())
			return
		case apperrors.SystemError:
			ctrl.logger.Error("System error in update payment status",
				slog.String("operation", appErr.Operation),
				slog.String("message", appErr.Message),
				slog.Any("error", appErr.Err))
			response.SendInternalServerError(c, ctrl.logger, appErr)
			return
		default:
			// For any other error types, treat as internal server error
			ctrl.logger.Error("Unknown error in update payment status", slog.String("error", err.Error()))
			response.SendInternalServerError(c, ctrl.logger, err)
			return
		}
	}

	paymentResponse := dto.PaymentResponse{
		PaymentID:     payment.ID,
		OrderID:       payment.OrderID,
		PaymentMethod: payment.PaymentMethod,
		TransactionID: payment.TransactionID,
		PaymentStatus: payment.PaymentStatus,
		PaymentDate:   payment.PaymentDate,
	}

	response.SendSuccess(c, http.StatusOK, "Payment status updated successfully", paymentResponse)
}

func (ctrl *paymentController) DeletePayment(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 32)
	if err != nil {
		response.SendBadRequestError(c, "Invalid payment ID")
		return
	}

	err = ctrl.paymentService.DeletePayment(uint(paymentID))
	if err != nil {
		// Handle different types of errors appropriately
		switch appErr := err.(type) {
		case apperrors.ValidationError:
			ctrl.logger.Info("Validation error in delete payment",
				slog.String("field", appErr.Field),
				slog.String("message", appErr.Message),
				slog.Any("value", appErr.Value))
			response.SendBadRequestError(c, appErr.Error())
			return
		case apperrors.BusinessRuleError:
			ctrl.logger.Info("Business rule error in delete payment",
				slog.String("rule", appErr.Rule),
				slog.String("message", appErr.Message))
			response.SendNotFoundError(c, appErr.Error())
			return
		case apperrors.SystemError:
			ctrl.logger.Error("System error in delete payment",
				slog.String("operation", appErr.Operation),
				slog.String("message", appErr.Message),
				slog.Any("error", appErr.Err))
			response.SendInternalServerError(c, ctrl.logger, appErr)
			return
		default:
			// For any other error types, treat as internal server error
			ctrl.logger.Error("Unknown error in delete payment", slog.String("error", err.Error()))
			response.SendInternalServerError(c, ctrl.logger, err)
			return
		}
	}

	response.SendSuccess(c, http.StatusOK, "Payment deleted successfully", nil)
}
