package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	JWTSecretKey      string `mapstructure:"JWT_SECRET_KEY"`

	DBMaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
	DBConnMaxIdleTime time.Duration `mapstructure:"DB_CONN_MAX_IDLE_TIME"`

	RedisAddr         string `mapstructure:"REDIS_ADDR"`
	RedisPassword     string `mapstructure:"REDIS_PASSWORD"`
	RedisDB           int    `mapstructure:"REDIS_DB"`
}

var AppConfig Config

func InitConfig(logger *slog.Logger) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")

	v.AutomaticEnv()

	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_MAX_OPEN_CONNS", 100)
	v.SetDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute)
	v.SetDefault("DB_CONN_MAX_IDLE_TIME", 1*time.Minute)

	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn(".env file not found, using default values and environment variables")
		} else {
			logger.Error("failed to read config file", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}

	if err := v.Unmarshal(&AppConfig); err != nil {
		logger.Error("failed to unmarshal config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Configuration loaded successfully")
}

