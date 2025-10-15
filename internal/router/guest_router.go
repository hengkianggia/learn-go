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

func SetupGuestRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	guestRepo := repository.NewGuestRepository(db)
	guestService := service.NewGuestService(guestRepo, logger)
	guestController := controller.NewGuestController(guestService, logger, db)

	guestRoutes := rg.Group("/guests")
	{
		guestRoutes.GET("/", guestController.GetAllGuests)
		guestRoutes.GET("/:slug", guestController.GetGuestBySlug)

		guestRoutes.Use(middleware.AuthMiddleware())
		{
			guestRoutes.POST("/", middleware.RoleMiddleware(model.Administrator, model.Organizer), guestController.CreateGuest)
		}
	}
}
