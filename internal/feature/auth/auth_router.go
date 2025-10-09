package auth

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	// Auto-migrate the User model, making the module self-contained
	db.AutoMigrate(&User{})

	userRepo := NewUserRepository(db)
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