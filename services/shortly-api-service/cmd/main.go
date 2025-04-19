package main

import (
	"os"
	"time"

	"shortly-api-service/config"
	"shortly-api-service/internal/database"
	"shortly-api-service/internal/routes"
	"shortly-api-service/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize slog logger
	utils.InitLogger()

	// Load environment variables
	if err := config.Init(); err != nil {
		utils.Log.Error("❌ Failed to load env", "error", err)
		os.Exit(1)
	}

	utils.Log.Info("✅ Environment variables loaded successfully")

	// Initialize Gin server
	server := gin.Default()

	// Database connection
	if err := database.ConnectDB(); err != nil {
		utils.Log.Error("❌ Failed to connect to database", "error", err)
		os.Exit(1)
	}

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
	routes.UrlRouter(api)
	routes.AuthRouter(api)
	routes.HealthRouter(api)
	routes.ProfileRouter(api)

	utils.Log.Info("🚀 Server is running", "port", config.AppConfig.PORT)

	if err := server.Run(":" + config.AppConfig.PORT); err != nil {
		utils.Log.Error("Failed to start server", "error", err)
	}

}
