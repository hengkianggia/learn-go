package speaker

import (
	"log/slog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupSpeakerRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	// speakerRepo := NewSpeakerRepository(db)
	// speakerService := NewSpeakerService(speakerRepo, logger)
	// speakerController := NewSpeakerController(speakerService, logger)

	// speakerRoutes := rg.Group("/speakers")
	// {
	// 	speakerRoutes.POST("/", speakerController.CreateSpeaker)
	// 	speakerRoutes.GET("/", speakerController.GetAllSpeakers)
	// 	speakerRoutes.GET("/:id", speakerController.GetSpeakerByID)
	// 	speakerRoutes.PUT("/:id", speakerController.UpdateSpeaker)
	// 	speakerRoutes.DELETE("/:id", speakerController.DeleteSpeaker)
	// }
}
