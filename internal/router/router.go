package router

import (
	"learn/internal/middleware"
	"learn/internal/model"
	"learn/internal/pkg/events"
	"learn/internal/repository"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		fields := []slog.Attr{
			slog.String("request_id", middleware.GetRequestID(c)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("route", c.FullPath()),
			slog.Int("status", c.Writer.Status()),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		if userCtx, ok := c.Get("user"); ok {
			if user, ok := userCtx.(model.User); ok {
				fields = append(fields, slog.Uint64("user_id", uint64(user.ID)))
			}
		}
		if len(c.Errors) > 0 {
			fields = append(fields, slog.String("error", c.Errors.String()))
		}

		logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "request", fields...)
	}
}

func SetupRouter(logger *slog.Logger, db *gorm.DB, eventBus *events.EventBus) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.RequestIDMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000" ||
				origin == "http://localhost:5173" ||
				origin == "http://localhost:5174" ||
				origin == "http://127.0.0.1:3000" ||
				origin == "http://127.0.0.1:5173"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(LoggerMiddleware(logger))

	gin.SetMode(gin.ReleaseMode)

	// Create repositories
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	eventRepo := repository.NewEventRepository(db)

	// Register event handlers
	RegisterEventHandlersWithRepos(eventBus, orderRepo, paymentRepo, ticketRepo, eventRepo, logger)

	// Buat grup utama untuk /api/v1
	apiV1 := r.Group("/api/v1")
	{
		SetupAuthRoutes(apiV1, db, logger)
		SetupVenueRoutes(apiV1, db, logger)
		SetupGuestRoutes(apiV1, db, logger)
		SetupEventRoutes(apiV1, db, logger)
		SetupOrderRoutes(apiV1, db, logger, eventBus)
		SetupPaymentRoutes(apiV1, db, logger, eventBus)
		SetupTicketRoutes(apiV1, db, logger)
		SetupAdminRoutes(apiV1, db, logger)
	}

	return r
}
