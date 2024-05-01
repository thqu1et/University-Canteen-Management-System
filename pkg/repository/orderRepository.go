package repository

import (
	"database/sql"
	"log"
	"postgresSQLProject/pkg/models"
	"time"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (repo *OrderRepository) CreateOrder(order *models.Order) error {
	query := `INSERT INTO orders (user_id, total_amount, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := repo.DB.QueryRow(query, order.UserID, order.TotalAmount, order.Status, time.Now()).Scan(&order.ID)
	if err != nil {
		log.Println("Failed to create order:", err)
		return err
	}
	return nil
}

func (repo *OrderRepository) GetOrdersByUserID(userID int) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT id, user_id, total_amount, status, created_at FROM orders WHERE user_id = $1`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		log.Println("Failed to retrieve orders:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
			log.Println("Failed to scan order:", err)
			continue
		}
		orders = append(orders, o)
	}
	return orders, nil
}
