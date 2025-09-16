package config

import (
	"fmt"
	"log"
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

func InitConfig() {
	v := viper.New()
	v.AddConfigPath(".") // Mencari file config di direktori saat ini
	v.SetConfigName(".env") // Nama file config tanpa ekstensi
	v.SetConfigType("env") // Tipe file config

	v.AutomaticEnv() // Membaca variabel lingkungan

	// Set default values for connection pool
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_MAX_OPEN_CONNS", 100)
	v.SetDefault("DB_CONN_MAX_LIFETIME", 5 * time.Minute)
	v.SetDefault("DB_CONN_MAX_IDLE_TIME", 1 * time.Minute)

	// Set default values for Redis
	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Error: .env file not found. Please create one.")
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	if err := v.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Error unmarshalling config: %s", err)
	}

	fmt.Println("Configuration loaded successfully!")
}
