package cmd

import (
	"learn/internal/config"
	"learn/internal/database"
	"learn/internal/migration"
	"learn/internal/pkg/logger"
	"log/slog"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `This command runs all pending database migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogger()

		config.InitConfig(log)
		db := database.InitDatabase(log)

		migrator := migration.NewMigrator(db, log)

		if err := migrator.Run(); err != nil {
			log.Error("Migration failed", slog.String("error", err.Error()))
			cmd.SilenceUsage = true
			return
		}

		log.Info("All migrations completed successfully")
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback [version]",
	Short: "Rollback migrations to a specific version",
	Long:  `This command rolls back migrations to a specific version.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogger()

		config.InitConfig(log)
		_ = database.InitDatabase(log) // Initialize database but don't use it yet

		var targetVersion string
		if len(args) > 0 {
			targetVersion = args[0]
		} else {
			log.Info("Rolling back to previous version")
			// Implement rollback logic here
			return
		}

		log.Info("Rollback command not fully implemented yet",
			slog.String("target_version", targetVersion))
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
}
