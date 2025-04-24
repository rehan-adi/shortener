package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT          string
	MONGO_URI     string
	MONGO_DB_NAME string
	REDIS_ADDR    string
}

var AppConfig Config

func Init() error {

	err := godotenv.Load()

	if err != nil {
		return err
	}

	AppConfig = Config{
		PORT:          GetEnvOrPanic("PORT"),
		MONGO_URI:     GetEnvOrPanic("MONGO_URI"),
		MONGO_DB_NAME: GetEnvOrPanic("MONGO_DB_NAME"),
		REDIS_ADDR:    GetEnvOrPanic("REDIS_ADDR"),
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
