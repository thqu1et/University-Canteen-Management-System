package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID    uint    `json:"order_id"`
	MenuItemID uint    `json:"menu_item_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price" gorm:"type:decimal(10,2)"`
}
