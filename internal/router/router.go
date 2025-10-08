package router

import (
	"learn/internal/auth"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logger.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", duration),
		)
	}
}

func SetupRouter(logger *slog.Logger) *gin.Engine {
	r := gin.Default()

	r.Use(LoggerMiddleware(logger))

	gin.SetMode(gin.ReleaseMode)

	// Buat grup utama untuk /api/v1
	apiV1 := r.Group("/api/v1")
	{
		// Daftarkan rute dari setiap modul di dalam grup ini
		auth.SetupAuthRoutes(apiV1, logger)
		// product.SetupProductRoutes(apiV1) // Contoh untuk modul produk
	}

	return r
}
