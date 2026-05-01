package router

import (
	"learn/internal/controller"
	"learn/internal/middleware"
	"learn/internal/pkg/events"
	"learn/internal/pkg/ratelimiter"
	"learn/internal/repository"
	"learn/internal/service"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPaymentRoutes(apiV1 *gin.RouterGroup, db *gorm.DB, logger *slog.Logger, eventBus *events.EventBus) {
	paymentRepository := repository.NewPaymentRepository(db)
	orderRepository := repository.NewOrderRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	eventRepository := repository.NewEventRepository(db)
	paymentService := service.NewPaymentService(paymentRepository, orderRepository, ticketRepository, eventRepository, logger, eventBus)
	paymentController := controller.NewPaymentController(paymentService, logger)

	paymentRouter := apiV1.Group("/payments")
	// Public Routes
	paymentRouter.POST("/midtrans-notification", ratelimiter.Limit("payment_notification", 60, time.Minute), paymentController.HandleNotification)

	paymentRouter.Use(middleware.AuthMiddleware())
	{
		paymentRouter.POST("/", ratelimiter.Limit("payment_create", 10, time.Minute), paymentController.CreatePayment)
		paymentRouter.GET("/:id", paymentController.GetPaymentByID)
		paymentRouter.GET("/order/:order_id", paymentController.GetPaymentByOrderID)
		paymentRouter.PUT("/:id", paymentController.UpdatePayment)
		paymentRouter.PATCH("/:id/status", ratelimiter.Limit("payment_status_update", 20, time.Minute), paymentController.UpdatePaymentStatus)
		paymentRouter.DELETE("/:id", paymentController.DeletePayment)
	}
}
