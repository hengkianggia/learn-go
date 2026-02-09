package controller

import (
	"learn/internal/dto"
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
		response.HandleAppError(c, err, ctrl.logger, "create payment")
		return
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
		response.HandleAppErrorWithNotFound(c, err, ctrl.logger, "get payment by ID")
		return
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
		response.HandleAppErrorWithNotFound(c, err, ctrl.logger, "get payment by order ID")
		return
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
		response.HandleAppErrorWithNotFound(c, err, ctrl.logger, "update payment")
		return
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
		response.HandleAppErrorWithNotFound(c, err, ctrl.logger, "update payment status")
		return
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
		response.HandleAppErrorWithNotFound(c, err, ctrl.logger, "delete payment")
		return
	}

	response.SendSuccess(c, http.StatusOK, "Payment deleted successfully", nil)
}
