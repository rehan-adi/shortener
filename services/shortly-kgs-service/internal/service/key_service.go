package service

import (
	"context"
	"shortly-kgs-service/config"
	"shortly-kgs-service/internal/constants"
	"shortly-kgs-service/internal/database"
	"shortly-kgs-service/internal/kgs"
	"shortly-kgs-service/internal/models"
	"shortly-kgs-service/internal/redis"
	"shortly-kgs-service/internal/utils"
	"shortly-proto/gen/key"
)

type KeyServiceServer struct {
	key.UnimplementedKeyServiceServer
}

func NewKeyServiceServer() key.KeyServiceServer {
	return &KeyServiceServer{}
}

func (s *KeyServiceServer) GetKey(ctx context.Context, req *key.Empty) (*key.KeyResponse, error) {

	queueLen, err := redis.RedisClient.LLen(context.Background(), constants.RedisQueueName).Result()

	if err != nil {
		return nil, err
	}

	if queueLen < constants.QueueLength {
		utils.Log.Info("Queue length is low, generating more keys")
		if err := kgs.GenerateKeys(1000); err != nil {
			return nil, err
		}
	}

	keyVal, err := redis.RedisClient.RPop(ctx, constants.RedisQueueName).Result()

	if err != nil {
		return nil, err
	}

	collection := database.MongoClient.Database(config.AppConfig.MONGO_DB_NAME).Collection("shortkeys")

	filter := models.ShortKey{Key: keyVal}

	update := models.ShortKey{
		Status: models.Used,
	}

	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return nil, err
	}

	return &key.KeyResponse{Key: keyVal}, nil
}
