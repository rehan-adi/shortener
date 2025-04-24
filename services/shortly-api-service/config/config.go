package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT        string
	DB_PORT     string
	DB_HOST     string
	DB_USER     string
	DB_NAME     string
	DB_PASSWORD string
	JWT_SECRET  string
	REDIS_ADDR  string
}

var AppConfig Config

func Init() error {

	err := godotenv.Load()

	if err != nil {
		return err
	}

	AppConfig = Config{
		PORT:        GetEnvOrPanic("PORT"),
		DB_PORT:     GetEnvOrPanic("DB_PORT"),
		DB_HOST:     GetEnvOrPanic("DB_HOST"),
		DB_USER:     GetEnvOrPanic("DB_USER"),
		DB_NAME:     GetEnvOrPanic("DB_NAME"),
		DB_PASSWORD: GetEnvOrPanic("DB_PASSWORD"),
		JWT_SECRET:  GetEnvOrPanic("JWT_SECRET"),
		REDIS_ADDR:  GetEnvOrPanic("REDIS_ADDR"),
	}

	return nil
}

func GetEnvOrPanic(key string) string {

	value := os.Getenv(key)

	if value == "" {
		panic(fmt.Sprintf("❌ Missing required environment variable: %s", key))
	}

	return value
}
