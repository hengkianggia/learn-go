package controller

import (
	"learn/internal/dto"
	apperrors "learn/internal/errors"
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
		// Handle different types of errors appropriately
		switch appErr := err.(type) {
		case apperrors.ValidationError:
			ctrl.logger.Info("Validation error in create order",
				slog.String("field", appErr.Field),
				slog.String("message", appErr.Message),
				slog.Any("value", appErr.Value))
			response.SendBadRequestError(c, appErr.Error())
			return
		case apperrors.BusinessRuleError:
			ctrl.logger.Info("Business rule error in create order",
				slog.String("rule", appErr.Rule),
				slog.String("message", appErr.Message))
			response.SendBadRequestError(c, appErr.Error())
			return
		case apperrors.SystemError:
			ctrl.logger.Error("System error in create order",
				slog.String("operation", appErr.Operation),
				slog.String("message", appErr.Message),
				slog.Any("error", appErr.Err))
			response.SendInternalServerError(c, ctrl.logger, appErr)
			return
		default:
			// For any other error types, treat as internal server error
			ctrl.logger.Error("Unknown error in create order", slog.String("error", err.Error()))
			response.SendInternalServerError(c, ctrl.logger, err)
			return
		}
	}

	response.SendSuccess(c, http.StatusCreated, "Order created successfully", dto.ToOrderResponse(*order))
}
