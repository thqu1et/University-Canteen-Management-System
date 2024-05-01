package repository

import (
	"database/sql"
	"log"
	"postgresSQLProject/pkg/models"
)

type MenuRepository struct {
	DB *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{DB: db}
}

// CreateMenuItem adds a new item to the menu
func (repo *MenuRepository) CreateMenuItem(menuItem *models.MenuItem) error {
	query := `INSERT INTO menu (name, description, price, available) VALUES ($1, $2, $3, $4) RETURNING id`
	err := repo.DB.QueryRow(query, menuItem.Name, menuItem.Description, menuItem.Price, menuItem.Available).Scan(&menuItem.ID)
	if err != nil {
		log.Println("Failed to insert menu item:", err)
		return err
	}
	return nil
}

// GetMenuItems retrieves all items on the menu
func (repo *MenuRepository) GetMenuItems() ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	query := "SELECT id, name, description, price, available FROM menu"
	rows, err := repo.DB.Query(query)
	if err != nil {
		log.Println("Failed to retrieve menu items:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mi models.MenuItem
		if err := rows.Scan(&mi.ID, &mi.Name, &mi.Description, &mi.Price, &mi.Available); err != nil {
			log.Println("Error scanning menu item:", err)
			continue
		}
		menuItems = append(menuItems, mi)
	}
	return menuItems, nil
}

// UpdateMenuItemAvailability updates the availability of a specific menu item
func (repo *MenuRepository) UpdateMenuItemAvailability(id int, available bool) error {
	query := `UPDATE menu SET available = $2 WHERE id = $1`
	_, err := repo.DB.Exec(query, id, available)
	if err != nil {
		log.Println("Failed to update menu item availability:", err)
		return err
	}
	return nil
}
