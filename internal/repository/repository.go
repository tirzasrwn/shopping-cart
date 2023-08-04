package repository

import (
	"database/sql"

	"github.com/tirzasrwn/shopping-cart/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	// user
	GetUserByEmail(email string) (*models.User, error)
	InsertUser(user *models.User) (int, error)
	CreateCartForNewUser(userID int) (int, error)
	// category
	GetCategories() ([]*models.Category, error)
	// product
	GetProductByCategory(id int) ([]*models.Product, error)
}
