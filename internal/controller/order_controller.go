package controller

import (
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/request"
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
	if !request.BindJSONOrError(c, &input, ctrl.logger, "create order") {
		return
	}

	user, exists := c.Get("user")
	if !exists {
		response.SendUnauthorizedError(c, "User not authenticated")
		return
	}

	order, err := ctrl.orderService.CreateOrder(input, user.(model.User).ID)
	if err != nil {
		// Handle different types of errors appropriately
		response.HandleAppError(c, err, ctrl.logger, "create order")
		return
	}

	response.SendSuccess(c, http.StatusCreated, "Order created successfully", dto.ToOrderResponse(*order))
}
