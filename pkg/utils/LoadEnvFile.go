package utils

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvFile() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
