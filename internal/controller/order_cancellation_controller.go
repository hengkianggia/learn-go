package controller

import (
	"learn/internal/model"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type orderCancellationController struct {
	orderCancellationService service.OrderCancellationService
	logger                   *slog.Logger
}

type OrderCancellationController interface {
	CancelOrder(c *gin.Context)
}

func NewOrderCancellationController(orderCancellationService service.OrderCancellationService, logger *slog.Logger) OrderCancellationController {
	return &orderCancellationController{
		orderCancellationService: orderCancellationService,
		logger:                   logger,
	}
}

func (ctrl *orderCancellationController) CancelOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		ctrl.logger.Warn("Invalid order ID format", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid order ID")
		return
	}

	user, exists := c.Get("user")
	if !exists {
		response.SendUnauthorizedError(c, "User not authenticated")
		return
	}

	userModel := user.(model.User)

	// Get reason from request body
	type CancelRequest struct {
		Reason string `json:"reason"`
	}

	var req CancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.Warn("Failed to bind JSON for cancel order", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid request format")
		return
	}

	if req.Reason == "" {
		req.Reason = "Manual cancellation by user"
	}

	err = ctrl.orderCancellationService.CancelOrder(uint(orderID), userModel.ID, req.Reason)
	if err != nil {
		ctrl.logger.Error("Failed to cancel order",
			slog.Uint64("order_id", uint64(orderID)),
			slog.String("error", err.Error()))
		response.SendBadRequestError(c, err.Error())
		return
	}

	response.SendSuccess(c, http.StatusOK, "Order cancelled successfully", nil)
}