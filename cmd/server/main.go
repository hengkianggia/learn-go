package main

import (
	"learn/internal/config"
	"learn/internal/database"
	"learn/internal/router"
)

func main() {
	// Inisialisasi konfigurasi dari .env
	config.InitConfig()

	// Inisialisasi database
	database.ConnectDatabase()

	// Setup router
	r := router.SetupRouter()

	// Jalankan server
	r.Run(":8080")
}