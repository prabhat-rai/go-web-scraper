package main

import (
	"echoApp/handler"
	"echoApp/middlewares"
	"echoApp/repositories"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func registerRoutes(e *echo.Echo, client *mongo.Client) *handler.Handler {
	// Connect to DB
	dbName := os.Getenv("DB_DATABASE")
	database := client.Database(dbName)

	// Initialize handler
	h := &handler.Handler{
		DB: database,
		UserRepository : &repositories.UserRepository{
			DB: database,
		},
	}

	// Register Routes
	e.GET("/login", h.LoginForm, middlewares.Guest)
	e.POST("/login", h.Login, middlewares.Guest)
	e.GET("/register", h.RegisterForm, middlewares.Guest)
	e.POST("/register", h.Register, middlewares.Guest)

	e.GET("/", h.Home, middlewares.Authenticated)
	e.GET("/list", h.List, middlewares.Authenticated)
	e.GET("/logout", h.Logout, middlewares.Authenticated)

	e.GET("/dev-test/verify-mongodb-queries", h.VerifyMongoDbQueries, middlewares.Authenticated)
	e.GET("/dev-test/review", h.FetchReview, middlewares.Authenticated)

	e.Static("/static", "public/static")

	fmt.Println("Registering Routes : DONE")

	return h
}
