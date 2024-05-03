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

func ProcessPayment(db *gorm.DB, payment Payment) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		if err := tx.Model(&Order{}).Where("id = ?", payment.OrderID).Update("status", "Paid").Error; err != nil {
			return err
		}
		return nil
	})
}
