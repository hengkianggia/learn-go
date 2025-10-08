package auth

import (
	"learn/internal/database"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes sekarang menerima *gin.RouterGroup.
func SetupAuthRoutes(rg *gin.RouterGroup) {
	// Inisialisasi semua layer untuk auth
	userRepo := NewUserRepository(database.DB)
	authService := NewAuthService(userRepo)
	authController := NewAuthController(authService)

	// Grup rute publik untuk auth, menjadi: /api/v1/auth
	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", authController.Logout)
	}

	// Grup rute yang dilindungi, menjadi: /api/v1/profile
	// Middleware hanya diterapkan di grup ini.
	protected := rg.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/profile", authController.Profile)
	}
}