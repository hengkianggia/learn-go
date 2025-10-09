package event

import (
	"learn/internal/feature/auth"
	"learn/internal/feature/venue"
	"log/slog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupEventRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	eventRepo := NewEventRepository(db)
	venueRepo := venue.NewVenueRepository(db)
	eventService := NewEventService(eventRepo, venueRepo, logger)
	eventController := NewEventController(eventService, logger, db)

	eventRoutes := rg.Group("/events")
	{
		eventRoutes.GET("/", eventController.GetAllEvents) // Public route
		
		// Authenticated routes
		authenticated := eventRoutes.Group("/")
		authenticated.Use(auth.AuthMiddleware())
		{
			authenticated.POST("/", auth.RoleMiddleware(auth.Administrator, auth.Organizer), eventController.CreateEvent)
		}
	}
}