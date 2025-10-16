package router

import (
	"learn/internal/controller"
	"learn/internal/middleware"
	"learn/internal/model"
	"learn/internal/repository"
	"learn/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrderRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, logger)
	orderController := controller.NewOrderController(orderService, logger)

	orderRoutes := rg.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware())
	{
		orderRoutes.POST("/", middleware.RoleMiddleware(model.Attendee), orderController.CreateOrder)
	}
}
