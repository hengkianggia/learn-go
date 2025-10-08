package auth

import (
	"learn/internal/database"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(rg *gin.RouterGroup, logger *slog.Logger) {
	userRepo := NewUserRepository(database.DB)
	authService := NewAuthService(userRepo, logger)
	authController := NewAuthController(authService, logger)

	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", authController.Logout)
	}

	protected := rg.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/profile", authController.Profile)
	}
}