package models

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	OrderID     uint      `json:"order_id"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2)"`
	Method      string    `json:"method"`
	ProcessedAt time.Time `json:"processed_at"`
}
