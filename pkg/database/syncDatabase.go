package database

import "postgresSQLProject/pkg/models"

func SyncDataBase() {
	DB.AutoMigrate(&models.User{}, &models.MenuItem{}, &models.OrderItem{})
}
