package main

import (
	"os"
	"shortly-api-service/config"
	"shortly-api-service/internal/database"
	"shortly-api-service/internal/models"
	"shortly-api-service/internal/utils"
)

func RunMigration() {

	utils.InitLogger()

	if err := config.Init(); err != nil {
		utils.Log.Error("❌ Failed to load env variables", "error", err)
		os.Exit(1)
	}

	if err := database.ConnectDB(); err != nil {
		utils.Log.Error("❌ Failed to connect to the database", "error", err)
		os.Exit(1)
	}

	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Url{},
		&models.Analytics{},
	)

	if err != nil {
		utils.Log.Error("❌ Migration failed", "error", err)
		os.Exit(1)
	}

	utils.Log.Info("✅ Database migration completed successfully")

}

func main() {
	RunMigration()
}
