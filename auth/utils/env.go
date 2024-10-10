package utils

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadSecrets() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file path provided. Using system environment variables.")
	}
}
