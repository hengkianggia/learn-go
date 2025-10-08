package cmd

import (
	"learn/internal/config"
	"learn/internal/database"
	"learn/internal/pkg/logger"
	"learn/internal/router"
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