package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id"`
	Total      float64     `json:"total" gorm:"type:decimal(10,2)"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
	Paid       bool        `json:"paid" gorm:"default:false"`
}
