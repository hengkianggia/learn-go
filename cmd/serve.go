package cmd

import (
	"learn/internal/auth"
	"learn/internal/config"
	"learn/internal/database"
	"learn/internal/router"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the API server",
	Long:  `This command starts the HTTP API server for the application.`, // Deskripsi panjang
	Run: func(cmd *cobra.Command, args []string) {
		// Inisialisasi semua library
		config.InitConfig()
		config.ConnectRedis()
		database.ConnectDatabase()

		// Auto-migrate models
		err := database.DB.AutoMigrate(&auth.User{})
		if err != nil {
			panic("Failed to migrate database! " + err.Error())
		}

		// Setup router
		r := router.SetupRouter()

		// Jalankan server
		r.Run(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}