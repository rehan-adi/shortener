package clients

import (
	"shortly-api-service/config"
	"shortly-api-service/internal/utils"
	"shortly-proto/gen/key"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var KGSClient key.KeyServiceClient

func InitKGSClient() {

	conn, err := grpc.NewClient(
		config.AppConfig.KGS_GRPC_ADDRESS,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		utils.Log.Error("Failed to connect to KGS gRPC service", "error", err)
		panic(err)
	}

	KGSClient = key.NewKeyServiceClient(conn)

	utils.Log.Info("Connected to KGS gRPC service")
}
