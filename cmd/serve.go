package cmd

import (
	"context"
	"errors"
	"learn/internal/config"
	"learn/internal/database"
	"learn/internal/model"
	"learn/internal/pkg/events"
	"learn/internal/pkg/logger"
	"learn/internal/pkg/queue"
	"learn/internal/router"
	seed "learn/internal/seed"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the API server",
	Long:  `This command starts the HTTP API server for the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Initialize Logger
		log := logger.NewLogger()

		// 2. Initialize Configurations and Connections
		config.InitConfig(log)
		db := database.InitDatabase(log)
		config.ConnectRedis(log)

		// 3. Initialize event bus with Redis for Streams
		eventBus := events.NewEventBus(config.Rdb, log)
		eventBus.Start()

		// 4. Initialize and start job queue
		jobQueue := queue.NewJobQueue(5, log)
		jobQueue.Start()

		// Environment-based migration
		env := config.AppConfig.AppEnv
		if env == "development" || env == "testing" {
			// Create enums
			createEnums(db)

			// Auto migrate only in development/testing
			db.AutoMigrate(&model.User{}, &model.Venue{}, &model.Guest{}, &model.Event{},
				&model.EventPrice{}, &model.EventGuest{}, &model.Order{}, &model.Ticket{},
				&model.Payment{}, &model.OrderLineItem{})

			log.Info("Auto-migration completed for development environment")
		} else {
			log.Info("Running in production mode - manual migrations expected")
		}

		// Seed the database (only in development)
		if env == "development" {
			seed.SeedUsers(db, log)
		}

		// 5. Setup Router with dependencies
		r := router.SetupRouter(log, db, eventBus)

		// 6. Run Server with graceful shutdown
		srv := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}

		go func() {
			log.Info("Starting server on port :8080")
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("failed to run server", slog.String("error", err.Error()))
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(quit)
		<-quit

		log.Info("Shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("server forced to shutdown", slog.String("error", err.Error()))
		} else {
			log.Info("Server shutdown completed")
		}

		jobQueue.Stop()
		eventBus.Stop()

		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Error("failed to close database connection", slog.String("error", err.Error()))
			}
		} else {
			log.Error("failed to get database connection", slog.String("error", err.Error()))
		}

		if config.Rdb != nil {
			if err := config.Rdb.Close(); err != nil {
				log.Error("failed to close Redis connection", slog.String("error", err.Error()))
			}
		}
	},
}

func createEnums(db *gorm.DB) {
	// Create the user_type enum
	db.Exec(`DO $$ BEGIN
		CREATE TYPE user_type AS ENUM ('organizer', 'attendee', 'administrator');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	// Create the event_status enum
	db.Exec(`DO $$ BEGIN
		CREATE TYPE event_status AS ENUM ('DRAFT', 'PUBLISHED', 'CANCELLED');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
