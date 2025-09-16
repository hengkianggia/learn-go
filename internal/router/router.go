package router

import (
	"learn/internal/auth"
	"learn/internal/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Inisialisasi semua layer
	// Menggunakan package auth yang baru
	userRepo := auth.NewUserRepository(database.DB)
	authService := auth.NewAuthService(userRepo)
	authController := auth.NewAuthController(authService)

	// Grup route
	public := r.Group("/auth")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
		public.POST("/logout", authController.Logout)
	}

	protected := r.Group("/api")
	// Menggunakan middleware dari package auth
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/profile", authController.Profile)
	}

	return r
}