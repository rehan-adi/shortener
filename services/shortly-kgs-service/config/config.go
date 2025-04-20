package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MONGO_URI string
	REDIS_ADDR string
}

var AppConfig Config

func Init() error {

	err := godotenv.Load()

	if err != nil {
		return err
	}

	AppConfig = Config{
		MONGO_URI: GetEnvOrPanic("MONGO_URI"),
		REDIS_ADDR: GetEnvOrPanic("REDIS_ADDR"),
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
