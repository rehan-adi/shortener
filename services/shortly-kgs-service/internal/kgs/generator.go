package kgs

import (
	"context"
	"time"

	"shortly-kgs-service/config"
	"shortly-kgs-service/internal/database"
	"shortly-kgs-service/internal/models"
	"shortly-kgs-service/internal/redis"
	"shortly-kgs-service/internal/utils"
)

const RedisQueue = "shortly-kgs-redis-queue"
const RedisCounter = "shortly-kgs-queue-counter"

func GenerateKeys(count int) error {

	for i := 0; i < count; i++ {

		// Get next number using Redis atomic counter
		id, err := redis.RedisClient.Incr(context.Background(), RedisCounter).Result()

		if err != nil {
			return err
		}

		key := utils.Base62Encode(id)

		// Store in MongoDB
		collection := database.MongoClient.Database(config.AppConfig.MONGO_DB_NAME).Collection("shortkeys")

		_, err = collection.InsertOne(context.Background(), models.ShortKey{
			Key:       key,
			Status:    "available",
			CreatedAt: time.Now(),
		})

		if err != nil {
			return err
		}

		// Push to Redis queue
		err = redis.RedisClient.LPush(context.Background(), RedisQueue, key).Err()

		if err != nil {
			return err
		}
	}

	utils.Log.Info("✅ Successfully generated and stored keys", "count", count)
	return nil
}
