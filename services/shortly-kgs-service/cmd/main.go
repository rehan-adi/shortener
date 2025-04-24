package main

import (
	"net"
	"os"

	"shortly-kgs-service/config"
	"shortly-kgs-service/internal/database"
	"shortly-kgs-service/internal/redis"
	"shortly-kgs-service/internal/service"
	"shortly-kgs-service/internal/utils"

	"google.golang.org/grpc"

	"shortly-proto/gen/key"
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

	listener, err := net.Listen("tcp", config.AppConfig.PORT)

	if err != nil {
		utils.Log.Error("❌ Failed to listen on port", "port", config.AppConfig.PORT, "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()

	key.RegisterKeyServiceServer(grpcServer, service.NewKeyServiceServer())

	utils.Log.Info("Shortly KGS Service is running...")

	if err := grpcServer.Serve(listener); err != nil {
		utils.Log.Error("❌ Failed to serve gRPC server", "error", err)
		os.Exit(1)
	}

}
