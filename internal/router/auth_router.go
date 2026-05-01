package router

import (
	"learn/internal/controller"
	"learn/internal/middleware"
	"learn/internal/pkg/ratelimiter"
	"learn/internal/repository"
	"learn/internal/service"
	"log/slog"
	"time"

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
		authRoutes.POST("/register", ratelimiter.Limit("auth_register", 5, time.Minute), authController.Register)
		authRoutes.POST("/verify-otp", ratelimiter.Limit("auth_verify_otp", 5, time.Minute), authController.VerifyOTP)
		authRoutes.POST("/login", ratelimiter.Limit("auth_login", 5, time.Minute), authController.Login)
		authRoutes.POST("/logout", authController.Logout)
	}

	protected := rg.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", authController.Profile)
	}
}
