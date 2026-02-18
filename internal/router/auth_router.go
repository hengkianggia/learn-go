package router

import (
	"learn/internal/controller"
	"learn/internal/middleware"
	"learn/internal/repository"
	"learn/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	userRepo := repository.NewUserRepository(db)
	emailService := service.NewEmailService(logger)
	authService := service.NewAuthService(userRepo, emailService, logger)
	authController := controller.NewAuthController(authService, logger)

	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/verify-otp", authController.VerifyOTP)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", authController.Logout)
	}

	protected := rg.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", authController.Profile)
	}
}
