package venue

import (
	"learn/internal/feature/auth"
	"log/slog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupVenueRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	venueRepo := NewVenueRepository(db)
	venueService := NewVenueService(venueRepo, logger)
	venueController := NewVenueController(venueService, logger, db)

	venueRoutes := rg.Group("/venues")
	{
		venueRoutes.GET("/", venueController.GetAllVenues) // Public route

		// Authenticated routes
		authenticated := venueRoutes.Group("/")
		authenticated.Use(auth.AuthMiddleware())
		{
			authenticated.POST("/", auth.RoleMiddleware(auth.Administrator, auth.Organizer), venueController.CreateVenue)
		}
	}
}
