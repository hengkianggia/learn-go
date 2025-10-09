package speaker

import (
	"learn/internal/feature/auth"
	"log/slog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupSpeakerRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	speakerRepo := NewSpeakerRepository(db)
	speakerService := NewSpeakerService(speakerRepo, logger)
	speakerController := NewSpeakerController(speakerService, logger)

	speakerRoutes := rg.Group("/speakers")
	speakerRoutes.Use(auth.AuthMiddleware())
	{
		speakerRoutes.POST("/", auth.RoleMiddleware(auth.Administrator, auth.Organizer), speakerController.CreateSpeaker)
	}
}