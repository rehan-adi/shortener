package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DB_PORT     string
	DB_HOST     string
	DB_USER     string
	DB_NAME     string
	DB_PASSWORD string
	JWT_SECRET  string
}

var AppConfig Config

func Init() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = Config{
		Port:        os.Getenv("PORT"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		JWT_SECRET:  os.Getenv("JWT_SECRET"),
	}

}
