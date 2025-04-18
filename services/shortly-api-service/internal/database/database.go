package database

import (
	"fmt"
	"time"

	"shortly-api-service/config"
	"shortly-api-service/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() error {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DB_HOST,
		config.AppConfig.DB_USER,
		config.AppConfig.DB_PASSWORD,
		config.AppConfig.DB_NAME,
		config.AppConfig.DB_PORT,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("❌ Failed to open database (GORM connection setup): %w", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		return fmt.Errorf("❌ Failed to get DB instance from GORM: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("❌ Failed to ping database: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	utils.Log.Info("✅ Database connected successfully")

	DB = db

	return nil

}
