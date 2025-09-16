package router

import (
	"learn/internal/controllers"
	"learn/internal/database"
	"learn/internal/middleware"
	"learn/internal/repositories"
	"learn/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Inisialisasi semua layer
	userRepo := repositories.NewUserRepository(database.DB)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// Grup route
	public := r.Group("/auth")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
		public.POST("/logout", authController.Logout)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", authController.Profile)
	}

	return r
}
