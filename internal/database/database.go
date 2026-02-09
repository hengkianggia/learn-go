package database

import (
	"fmt"
	"learn/internal/config"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(logger *slog.Logger) *gorm.DB {
	host := config.AppConfig.DBHost
	port := config.AppConfig.DBPort
	user := config.AppConfig.DBUser
	password := config.AppConfig.DBPassword
	dbname := config.AppConfig.DBName

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	sqlDB, err := database.DB()
	if err != nil {
		logger.Error("failed to get underlying sql.DB object", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(config.AppConfig.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.AppConfig.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.AppConfig.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.AppConfig.DBConnMaxIdleTime)

	DB = database
	logger.Info("Successfully connected to database")
	return DB
}
