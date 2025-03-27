package main

import (
	"shortly-api-service/config"
	"shortly-api-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize logger
	utils.InitLogger()

	// Load environment variables
	config.Init()
	utils.Log.Info("✅ Environment variables loaded successfully")

	// Initialize Gin server
	server := gin.Default()

    // api = server.Group("/api/v1")


	utils.Log.Infof("🚀 Server running on port %s", config.AppConfig.Port)

	if err := server.Run(":" + config.AppConfig.Port); err != nil {
		utils.Log.Fatalf("❌ Failed to start server: %v", err)
	}

}
