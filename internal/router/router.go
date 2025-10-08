package router

import (
	"learn/internal/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	// Buat grup utama untuk /api/v1
	apiV1 := r.Group("/api/v1")
	{
		// Daftarkan rute dari setiap modul di dalam grup ini
		auth.SetupAuthRoutes(apiV1)
		// product.SetupProductRoutes(apiV1) // Contoh untuk modul produk
	}

	return r
}
