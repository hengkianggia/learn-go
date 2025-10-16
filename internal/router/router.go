package router

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func SetupRouter(logger *slog.Logger, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(LoggerMiddleware(logger))

	gin.SetMode(gin.ReleaseMode)

	// Buat grup utama untuk /api/v1
	apiV1 := r.Group("/api/v1")
	{
		SetupAuthRoutes(apiV1, db, logger)
		SetupVenueRoutes(apiV1, db, logger)
		SetupGuestRoutes(apiV1, db, logger)
		SetupEventRoutes(apiV1, db, logger)
		SetupOrderRoutes(apiV1, db, logger)
	}

	return r
}
