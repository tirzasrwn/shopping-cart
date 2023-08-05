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
	// order
	InsertOrder(cartID int, productID int, quantity int) (int, error)
	CheckOrderExist(cartID int, productID int) (bool, error)
	UpdateQuantity(cartID int, productID int, quantityToAdd int) (int, error)
	GetOrderByUserEmail(email string) ([]*models.Product, error)
}
