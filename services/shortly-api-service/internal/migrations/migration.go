package main

import (
	"shortly-api-service/config"
	"shortly-api-service/internal/database"
	"shortly-api-service/internal/models"
	"shortly-api-service/internal/utils"
)

func RunMigration() {

	utils.InitLogger()

	config.Init()

	database.ConnectDB()

	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Url{},
		&models.Analytics{},
	)

	if err != nil {
		utils.Log.Fatalf("❌ Migration failed: %v", err)
	}

	utils.Log.Info("✅ Database migration completed successfully")

}

func main() {
	RunMigration()
}
