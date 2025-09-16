package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     AppConfig.RedisAddr,
		Password: AppConfig.RedisPassword,
		DB:       AppConfig.RedisDB,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis! " + err.Error())
	}

	fmt.Println("Successfully connected to Redis!")
}
