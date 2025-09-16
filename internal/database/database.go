package database

import (
	"fmt"
	"learn/internal/config"
	"learn/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := config.AppConfig.DBHost
	port := config.AppConfig.DBPort
	user := config.AppConfig.DBUser
	password := config.AppConfig.DBPassword
	dbname := config.AppConfig.DBName

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database! " + err.Error())
	}

	err = database.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to migrate database! " + err.Error())
	}

	DB = database
	fmt.Println("Successfully connected to database!")
}
