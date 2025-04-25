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

	for i := 0; i < count; i++ {

		key, err := utils.GenerateRandomKey(6)

		if err != nil {
			return err
		}

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
		err = redis.RedisClient.LPush(context.Background(), constants.RedisQueueName, key).Err()

		if err != nil {
			return err
		}
	}

	utils.Log.Info("✅ Successfully generated and stored keys", "count", count)
	return nil
}
