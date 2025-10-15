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

func SetupEventRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	eventRepo := repository.NewEventRepository(db)
	venueRepo := repository.NewVenueRepository(db)
	guestRepo := repository.NewGuestRepository(db)
	eventService := service.NewEventService(eventRepo, venueRepo, guestRepo, logger)
	eventController := controller.NewEventController(eventService, logger, db)

	eventRoutes := rg.Group("/events")
	{
		eventRoutes.GET("/", eventController.GetAllEvents) // Public route
		eventRoutes.GET("/:slug", eventController.GetEventBySlug)

		// Authenticated routes
		authenticated := eventRoutes.Group("/")
		authenticated.Use(middleware.AuthMiddleware())
		{
			authenticated.POST("/", middleware.RoleMiddleware(model.Administrator, model.Organizer), eventController.CreateEvent)
		}
	}
}
