package cmd

import (
	"learn/internal/config"
	"learn/internal/database"
	"learn/internal/model"
	"learn/internal/pkg/events"
	"learn/internal/pkg/logger"
	"learn/internal/pkg/queue"
	"learn/internal/router"
	seed "learn/internal/seed"
	"log/slog"
	"os"

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

		// 3. Initialize event bus
		eventBus := events.NewEventBus()

		// 4. Initialize and start job queue
		jobQueue := queue.NewJobQueue(5, log)
		jobQueue.Start()
		defer jobQueue.Stop()

		// Environment-based migration
		env := os.Getenv("APP_ENV")
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

		// 6. Run Server
		log.Info("Starting server on port :8080")
		if err := r.Run(":8080"); err != nil {
			log.Error("failed to run server", slog.String("error", err.Error()))
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
