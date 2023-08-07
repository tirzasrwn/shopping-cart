package handlers

import (
	"database/sql"
	"log"

	"github.com/tirzasrwn/shopping-cart/configs"
	"github.com/tirzasrwn/shopping-cart/internal/models"
	"github.com/tirzasrwn/shopping-cart/internal/repository"
	"github.com/tirzasrwn/shopping-cart/internal/repository/dbrepo"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB(dsn string) (*sql.DB, error) {
	connection, err := openDB(dsn)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres!")
	return connection, nil
}

func InitializeHandler(app *configs.Config) (appDB repository.DatabaseRepo, err error) {
	db, err := connectToDB(app.DSN)
	if err != nil {
		log.Println("[INIT] failed connecting to PostgreSQL")
		return
	}
	Handlers = &module{
		db: &dbEntity{
			conn:   db,
			dbrepo: &dbrepo.PostgresDBRepo{DB: db},
		},
	}
	return &dbrepo.PostgresDBRepo{DB: db}, nil
}
