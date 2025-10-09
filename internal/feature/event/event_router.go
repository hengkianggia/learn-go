package event

import (
	"learn/internal/feature/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupEventRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	eventRepo := NewEventRepository(db)
	eventService := NewEventService(eventRepo, logger)
	eventController := NewEventController(eventService, logger)

	eventRoutes := rg.Group("/events")
	eventRoutes.Use(auth.AuthMiddleware()) // All event routes need authentication
	{
		eventRoutes.POST("/", auth.RoleMiddleware(auth.Administrator, auth.Organizer), eventController.CreateEvent)
	}
}
