package database

import "postgresSQLProject/pkg/models"

func SyncDataBase() {
	DB.AutoMigrate(&models.User{})
}
