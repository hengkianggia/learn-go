package event

import (
	"learn/internal/feature/auth"
	"learn/internal/feature/guest"
	"learn/internal/feature/venue"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupEventRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	eventRepo := NewEventRepository(db)
	venueRepo := venue.NewVenueRepository(db)
	guestRepo := guest.NewGuestRepository(db)
	eventService := NewEventService(eventRepo, venueRepo, guestRepo, logger)
	eventController := NewEventController(eventService, logger, db)

	eventRoutes := rg.Group("/events")
	{
		eventRoutes.GET("/", eventController.GetAllEvents) // Public route
		eventRoutes.GET("/:slug", eventController.GetEventBySlug)

		// Authenticated routes
		authenticated := eventRoutes.Group("/")
		authenticated.Use(auth.AuthMiddleware())
		{
			authenticated.POST("/", auth.RoleMiddleware(auth.Administrator, auth.Organizer), eventController.CreateEvent)
		}
	}
}