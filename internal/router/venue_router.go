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

func SetupVenueRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	venueRepo := repository.NewVenueRepository(db)
	venueService := service.NewVenueService(venueRepo, logger)
	venueController := controller.NewVenueController(venueService, logger, db)

	venueRoutes := rg.Group("/venues")
	{
		venueRoutes.GET("/", venueController.GetAllVenues)
		venueRoutes.GET("/:slug", venueController.GetVenueBySlug)

		// Authenticated routes
		authenticated := venueRoutes.Group("/")
		authenticated.Use(middleware.AuthMiddleware())
		{
			authenticated.POST("/", middleware.RoleMiddleware(model.Administrator, model.Organizer), venueController.CreateVenue)
		}
	}
}
