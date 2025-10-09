package router

import (
	"learn/internal/feature/auth"
	"learn/internal/feature/event"
	"learn/internal/feature/speaker"
	"learn/internal/feature/venue"
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
		auth.SetupAuthRoutes(apiV1, db, logger)
		venue.SetupVenueRoutes(apiV1, db, logger)
		speaker.SetupSpeakerRoutes(apiV1, db, logger)
		event.SetupEventRoutes(apiV1, db, logger)
	}

	return r
}
