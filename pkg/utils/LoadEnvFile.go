package utils

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvFile() {
	err := godotenv.Load("/Users/askarabylkhaiyrov/Desktop/GoLang/postgresSQLProject/pkg/database/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
