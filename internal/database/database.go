package database

import (
	"fmt"
	"learn/internal/config"

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

	sqlDB, err := database.DB()
	if err != nil {
		panic("Failed to get underlying DB object! " + err.Error())
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(config.AppConfig.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.AppConfig.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.AppConfig.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.AppConfig.DBConnMaxIdleTime)

	DB = database
	fmt.Println("Successfully connected to database!")
}