package database

import (
	"context"
	"time"

	"shortly-kgs-service/config"
	"shortly-kgs-service/internal/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectDB() error {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(config.AppConfig.MONGO_URI)

	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		utils.Log.Error("❌ Failed to connect to MongoDB", "error", err)
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		utils.Log.Error("❌ Could not ping MongoDB", "error", err)
		return err
	}

	utils.Log.Info("✅ Connected to MongoDB")
	MongoClient = client

	return nil
}

func CloseMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := MongoClient.Disconnect(ctx); err != nil {
			utils.Log.Warn("Error disconnecting MongoDB", "error", err)
		} else {
			utils.Log.Info("MongoDB connection closed")
		}
	}
}
