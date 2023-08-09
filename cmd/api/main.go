package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/tirzasrwn/shopping-cart/configs"
	_ "github.com/tirzasrwn/shopping-cart/docs"
	"github.com/tirzasrwn/shopping-cart/internal/controllers/middleware"
	"github.com/tirzasrwn/shopping-cart/internal/handlers"
	"github.com/tirzasrwn/shopping-cart/internal/logs"
	"github.com/tirzasrwn/shopping-cart/internal/routes"
)

func init() {
	configs.InitializeAppConfig()
	if !configs.AppConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

//	@title			Shopping Cart
//	@version		1.2.0
//	@description	This page is API documentation for API related to shopping-cart
//	@schemes		http
//	@host			localhost:4000
//	@BasePath		/v1
//	@contact.name	tirzasrwn
//	@contact.email	tirzasrwn@gmail.com

//	@securityDefinitions.apikey	UserAuth
//	@in							header
//	@name						Authorization

func main() {
	app := configs.AppConfig
	fmt.Println(app)

	db := connectToDB(app.DSN)
	if db == nil {
		log.Panic("can't connect to database")
	}
	defer db.Close()
	app.DB = db

	mongoClient, err := connectToMongo(app.MongoUrl)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("success connect to mongo", &mongoClient)
	defer func(client *mongo.Client) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}(mongoClient)
	err = logs.New(mongoClient)
	if err != nil {
		log.Panic(err)
		return
	}

	err = handlers.InitializeHandler(&app)
	if err != nil {
		log.Panic(err)
		return
	}

	middleware.InitializeAuthenticationMiddleware(&app)
	if err != nil {
		log.Panic(err)
		return
	}

	routes := routes.InitializeRouter()
	routes.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", app.Port),
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func connectToDB(dsn string) *sql.DB {
	counts := 0
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready ...")
		} else {
			log.Println("connected to database")
			return connection
		}
		if counts > 10 {
			return nil
		}
		log.Print("backing off for 1 second")
		time.Sleep(time.Second)
		counts++
		continue
	}
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

func connectToMongo(mongoURL string) (*mongo.Client, error) {
	// Create connection options.
	clientOption := options.Client().ApplyURI(mongoURL)

	// Connect.
	log.Println("Connecting to mongo ...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return c, nil
}
