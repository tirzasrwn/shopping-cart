package handlers

import (
	"database/sql"

	"github.com/tirzasrwn/shopping-cart/configs"
	"github.com/tirzasrwn/shopping-cart/internal/models"
	"github.com/tirzasrwn/shopping-cart/internal/repository"
	"github.com/tirzasrwn/shopping-cart/internal/repository/dbrepo"
)

type HandlerFunc interface {
	// user
	GetUserByEmail(email string) (*models.User, error)
	InsertUser(user *models.User) (int, error)
	GetUserOrder(email string) ([]*models.ProductOrder, error)
	GetUserPayment(email string) ([]*models.ProductPayment, error)
	GetUserCartByEmail(email string) (int, error)
	// category
	GetCategory() ([]*models.Category, error)
	// product
	GetProducts() ([]*models.Product, error)
	GetProductByCategory(id int) ([]*models.Product, error)
	// order
	InsertOrder(cartID int, productID int, quantity int) (int, error)
	DeleteOrder(orderID int) error
	CheckoutOrder(money float64, email string) (float64, error)
}

var Handlers HandlerFunc

type module struct {
	db *dbEntity
}

type dbEntity struct {
	conn   *sql.DB
	dbrepo repository.DatabaseRepo
}

func InitializeHandler(app *configs.Config) (err error) {
	Handlers = &module{
		db: &dbEntity{
			conn: app.DB,
			dbrepo: &dbrepo.PostgresDBRepo{
				DB: app.DB,
			},
		},
	}
	return nil
}
