package guest

import (
	"learn/internal/feature/auth"
	"log/slog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupGuestRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	guestRepo := NewGuestRepository(db)
	guestService := NewGuestService(guestRepo, logger)
	guestController := NewGuestController(guestService, logger)

	guestRoutes := rg.Group("/guests")
	guestRoutes.Use(auth.AuthMiddleware())
	{
		guestRoutes.POST("/", auth.RoleMiddleware(auth.Administrator, auth.Organizer), guestController.CreateGuest)
	}
}
