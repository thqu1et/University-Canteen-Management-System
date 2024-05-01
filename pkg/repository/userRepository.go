package repository

import (
	"database/sql"
	"log"
	"postgresSQLProject/pkg/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	err := repo.DB.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		log.Println("Failed to insert new user:", err)
		return err
	}
	return nil
}

func (repo *UserRepository) GetUser(email string) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, first_name, last_name, email FROM users WHERE email = $1`
	err := repo.DB.QueryRow(query, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		log.Println("Failed to get user:", err)
		return nil, err
	}
	return u, nil
}

func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	query := "SELECT id, first_name, last_name, email FROM users"
	rows, err := repo.DB.Query(query)
	if err != nil {
		log.Println("Failed to retrieve users:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err != nil {
			log.Println("Error scanning user:", err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}
