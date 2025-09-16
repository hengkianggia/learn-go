package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

var AppConfig Config

func InitConfig() {
	v := viper.New()
	v.AddConfigPath(".") // Mencari file config di direktori saat ini
	v.SetConfigName(".env") // Nama file config tanpa ekstensi
	v.SetConfigType("env") // Tipe file config

	v.AutomaticEnv() // Membaca variabel lingkungan

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
