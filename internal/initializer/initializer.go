package initializer

import (
	"learn/internal/config"
	"learn/internal/database"
)

func InitApp() {
	config.InitConfig()
	config.ConnectRedis()
	database.ConnectDatabase()
}
