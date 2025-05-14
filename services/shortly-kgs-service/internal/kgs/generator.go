package kgs

import (
	"context"
	"time"

	"shortly-kgs-service/config"
	"shortly-kgs-service/internal/constants"
	"shortly-kgs-service/internal/database"
	"shortly-kgs-service/internal/models"
	"shortly-kgs-service/internal/redis"
	"shortly-kgs-service/internal/utils"
)

func GenerateKeys(count int) error {

	ctx := context.Background()
	collection := database.MongoClient.Database(config.AppConfig.MONGO_DB_NAME).Collection("shortkeys")

	var keys []interface{}
	var redisKeys []string

	for i := 0; i < count; i++ {

		key, err := utils.GenerateRandomKey(6)

		if err != nil {
			return err
		}

		keys = append(keys, models.ShortKey{
			Key:       key,
			Status:    "available",
			CreatedAt: time.Now(),
		})

		redisKeys = append(redisKeys, key)

	}

	if len(keys) > 0 {
		_, err := collection.InsertMany(ctx, keys)

		if err != nil {
			return err
		}
	}

	if len(redisKeys) > 0 {
		err := redis.RedisClient.LPush(ctx, constants.RedisQueueName, redisKeys).Err()

		if err != nil {
			return err
		}
	}

	utils.Log.Info("✅ Successfully generated and stored keys", "count", count)
	return nil
}
