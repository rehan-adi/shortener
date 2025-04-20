package main

import (
	"os"

	"shortly-kgs-service/config"
	"shortly-kgs-service/internal/database"
	"shortly-kgs-service/internal/redis"
	"shortly-kgs-service/internal/utils"
)

func main() {

	utils.InitLogger()

	if err := config.Init(); err != nil {
		utils.Log.Error("❌ Failed to load env", "error", err)
		os.Exit(1)
	}

	utils.Log.Info("✅ Environment variables loaded successfully")

	if err := database.ConnectDB(); err != nil {
		utils.Log.Error("MongoDB connection failed", "error", err)
		os.Exit(1)
	}

	defer database.CloseMongoDB()

	if err := redis.ConnectRedis(); err != nil {
		utils.Log.Error("❌ Failed to connect to Redis", "error", err)
		os.Exit(1)
	}

	defer redis.RedisClient.Close()

	utils.Log.Info("Shortly KGS Service is running...")

	select {}
}
