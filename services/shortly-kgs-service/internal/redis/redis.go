package redis

import (
	"context"

	"shortly-kgs-service/config"
	"shortly-kgs-service/internal/utils"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() error {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.REDIS_ADDR,
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping(context.Background()).Result()

	if err != nil {
		utils.Log.Error("❌ Redis connection failed", "error", err)
		return err
	}

	utils.Log.Info("✅ Connected to Redis")

	return nil

}
