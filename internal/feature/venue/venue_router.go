package venue

import (
	"log/slog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupVenueRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	// venueRepo := NewVenueRepository(db)
	// venueService := NewVenueService(venueRepo, logger)
	// venueController := NewVenueController(venueService, logger)

	// venueRoutes := rg.Group("/venues")
	// {
	// 	venueRoutes.POST("/", venueController.CreateVenue)
	// 	venueRoutes.GET("/", venueController.GetAllVenues)
	// 	venueRoutes.GET("/:id", venueController.GetVenueByID)
	// 	venueRoutes.PUT("/:id", venueController.UpdateVenue)
	// 	venueRoutes.DELETE("/:id", venueController.DeleteVenue)
	// }
}
