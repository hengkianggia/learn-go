package cmd

import (
	"learn/internal/config"
	"learn/internal/database"
	"learn/internal/model"
	"learn/internal/pkg/logger"
	"learn/internal/router"
	seed "learn/internal/seed"
	"log/slog"

	"github.com/spf13/cobra"
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

		// Drop table for development
		// db.Migrator().DropTable(&model.User{}, &model.Venue{}, &model.Guest{}, &model.Event{}, &model.EventPrice{}, &model.EventGuest{})

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

		// Migrate the schema
		db.AutoMigrate(&model.User{}, &model.Venue{}, &model.Guest{}, &model.Event{}, &model.EventPrice{}, &model.EventGuest{})

		// Seed the database
		seed.SeedUsers(db, log)

		// 3. Setup Router with dependencies
		r := router.SetupRouter(log, db)

		// 4. Run Server
		log.Info("Starting server on port :8080")
		if err := r.Run(":8080"); err != nil {
			log.Error("failed to run server", slog.String("error", err.Error()))
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
