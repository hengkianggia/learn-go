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

func SetupAdminRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	userRepo := repository.NewUserRepository(db)
	emailService := service.NewEmailService(logger)
	adminService := service.NewAdminService(userRepo, emailService, logger)
	adminController := controller.NewAdminController(adminService, logger, db)

	adminRoutes := rg.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware(model.Administrator))
	{
		adminRoutes.POST("/users/approve", adminController.ApproveUser)
		adminRoutes.POST("/users/reject", adminController.RejectUser)
		adminRoutes.POST("/users/block", adminController.BlockUser)
		adminRoutes.POST("/users/unblock", adminController.UnblockUser)
		adminRoutes.POST("/users/delete", adminController.DeleteUser)
		adminRoutes.GET("/users", adminController.ListUsers)
	}
}
