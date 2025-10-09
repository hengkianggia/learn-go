package event

import (
	"log/slog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupEventRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	// eventRepo := NewEventRepository(db)
	// eventService := NewEventService(eventRepo, logger)
	// eventController := NewEventController(eventService, logger)

	// eventRoutes := rg.Group("/events")
	// {
	// 	eventRoutes.POST("/", eventController.CreateEvent)
	// 	eventRoutes.GET("/", eventController.GetAllEvents)
	// 	eventRoutes.GET("/:id", eventController.GetEventByID)
	// 	eventRoutes.PUT("/:id", eventController.UpdateEvent)
	// 	eventRoutes.DELETE("/:id", eventController.DeleteEvent)
	// }
}
