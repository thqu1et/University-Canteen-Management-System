package database

import (
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDB() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, username, password, dbName, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object from GORM: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected!")
}
