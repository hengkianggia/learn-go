package cmd

import (
	"learn/internal/initializer"
	"learn/internal/router"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the API server",
	Long:  `This command starts the HTTP API server for the application.`, // Deskripsi panjang
	Run: func(cmd *cobra.Command, args []string) {
		// Inisialisasi semua library
		initializer.InitApp()

		// Setup router
		r := router.SetupRouter()

		// Jalankan server
		r.Run(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
