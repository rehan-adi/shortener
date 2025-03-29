package main

import (
	"shortly-api-service/config"
	"shortly-api-service/internal/database"
	"shortly-api-service/internal/routes"
	"shortly-api-service/internal/utils"
	"time"

	"github.com/gin-contrib/cors"
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

	// Database connection
	database.ConnectDB()

	// Middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := server.Group("/api/v1")

	// Routes
	routes.HealthRouter(api)
	routes.AuthRouter(api)

	utils.Log.Infof("🚀 Server running on port %s", config.AppConfig.Port)

	if err := server.Run(":" + config.AppConfig.Port); err != nil {
		utils.Log.Fatalf("❌ Failed to start server: %v", err)
	}

}
