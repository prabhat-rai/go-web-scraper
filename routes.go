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
		AppReviewRepository: &repositories.AppReviewRepository{
			DB: database,
		},
		KeywordRepository: &repositories.KeywordRepository{
			DB: database,
		},
		KeywordGroupRepository: &repositories.KeywordGroupRepository{
			DB: database,
		},
	}

	// Auth Routes
	e.GET("/login", h.LoginForm, middlewares.Guest)
	e.POST("/login", h.Login, middlewares.Guest)
	e.GET("/register", h.RegisterForm, middlewares.Guest)
	e.POST("/register", h.Register, middlewares.Guest)

	// Authenticated Routes
	e.GET("/", h.Home, middlewares.Authenticated)
	e.GET("/logout", h.Logout, middlewares.Authenticated)

	// Listing Routes
	e.GET("/apps", h.AppsList, middlewares.Authenticated)
	e.GET("/reviews", h.ListReviews, middlewares.Authenticated)
	e.GET("/keywords", h.ListKeywords, middlewares.Authenticated)
	e.GET("/keyword-groups", h.ListKeywordGroups, middlewares.Authenticated)

	// AJAX listing Routes
	e.GET("/ajax/reviews-list", h.RetrieveReviews, middlewares.Authenticated)
	e.GET("/ajax/keywords-list", h.RetrieveKeywords, middlewares.Authenticated)
	e.GET("/ajax/keyword-groups-list", h.RetrieveKeywordGroups, middlewares.Authenticated)

	// Dev Test Routes
	e.GET("/dev-test/verify-mongodb-queries", h.VerifyMongoDbQueries, middlewares.Authenticated)
	e.GET("/dev-test/review", h.FetchReview, middlewares.Authenticated)

	// File server
	e.Static("/static", "public/static")

	fmt.Println("Registering Routes : DONE")

	return h
}
