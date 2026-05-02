package router

import (
	"learn/internal/controller"
	"learn/internal/middleware"
	"learn/internal/model"
	"learn/internal/pkg/ratelimiter"
	"learn/internal/repository"
	"learn/internal/service"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTicketRoutes(apiV1 *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	ticketRepository := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepository, logger)
	ticketController := controller.NewTicketController(ticketService, logger)

	ticketRoutes := apiV1.Group("/tickets")
	ticketRoutes.Use(middleware.AuthMiddleware())
	{
		ticketRoutes.POST(
			"/check-in",
			middleware.RoleMiddleware(model.Administrator, model.Organizer),
			ratelimiter.Limit("ticket_checkin", 120, time.Minute),
			ticketController.CheckInTicket,
		)
	}
}
