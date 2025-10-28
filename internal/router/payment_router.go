package router

import (
	"learn/internal/controller"
	"learn/internal/middleware"
	"learn/internal/repository"
	"learn/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPaymentRoutes(apiV1 *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	paymentRepository := repository.NewPaymentRepository(db)
	orderRepository := repository.NewOrderRepository(db)
	paymentService := service.NewPaymentService(paymentRepository, orderRepository, logger)
	paymentController := controller.NewPaymentController(paymentService, logger)

	paymentRouter := apiV1.Group("/payments")
	paymentRouter.Use(middleware.AuthMiddleware())
	{
		paymentRouter.POST("/", paymentController.CreatePayment)
		paymentRouter.GET("/:id", paymentController.GetPaymentByID)
		paymentRouter.GET("/order/:order_id", paymentController.GetPaymentByOrderID)
		paymentRouter.PUT("/:id", paymentController.UpdatePayment)
		paymentRouter.PATCH("/:id/status", paymentController.UpdatePaymentStatus)
		paymentRouter.DELETE("/:id", paymentController.DeletePayment)
	}
}
