package repository

import (
	"database/sql"

	"github.com/tirzasrwn/shopping-cart/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	GetUserByEmail(email string) (*models.User, error)
	InsertUser(user *models.User) error
	// GetUserByID(id int) (*models.User, error)
}
