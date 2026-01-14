package router

import (
	"learn/internal/controller"
	"learn/internal/middleware"
	"learn/internal/model"
	"learn/internal/pkg/events"
	"learn/internal/repository"
	"learn/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrderRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger, eventBus *events.EventBus) {
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, logger, eventBus)
	orderController := controller.NewOrderController(orderService, logger)

	// Order cancellation service and controller
	orderCancellationService := service.NewOrderCancellationService(orderRepo, logger, eventBus)
	orderCancellationController := controller.NewOrderCancellationController(orderCancellationService, logger)

	orderRoutes := rg.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware())
	{
		orderRoutes.POST("/", middleware.RoleMiddleware(model.Attendee), orderController.CreateOrder)
		orderRoutes.DELETE("/:id", middleware.RoleMiddleware(model.Attendee), orderCancellationController.CancelOrder)
	}
}
