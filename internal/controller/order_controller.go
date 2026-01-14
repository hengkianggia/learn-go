package controller

import (
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type orderController struct {
	orderService service.OrderService
	logger       *slog.Logger
}
type OrderController interface {
	CreateOrder(c *gin.Context)
}

func NewOrderController(orderService service.OrderService, logger *slog.Logger) OrderController {
	return &orderController{orderService: orderService, logger: logger}
}

func (ctrl *orderController) CreateOrder(c *gin.Context) {
	var input dto.NewOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Warn("failed to bind JSON for create order", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return
	}

	user, exists := c.Get("user")
	if !exists {
		response.SendUnauthorizedError(c, "User not authenticated")
		return
	}

	order, err := ctrl.orderService.CreateOrder(input, user.(model.User).ID)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, "Order created successfully", dto.ToOrderResponse(*order))
}
