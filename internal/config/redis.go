package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var Ctx = context.Background()

func ConnectRedis(logger *slog.Logger) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     AppConfig.RedisAddr,
		Password: AppConfig.RedisPassword,
		DB:       AppConfig.RedisDB,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		logger.Error("failed to connect to Redis", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Successfully connected to Redis")
}
