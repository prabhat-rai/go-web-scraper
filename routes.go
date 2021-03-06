package main

import (
	"echoApp/handler"
	"echoApp/middlewares"
	"echoApp/repositories"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func registerRoutes(e *echo.Echo, client *mongo.Client,dbName string) *handler.Handler {
	// Connect to DB
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
		AppRepository: &repositories.AppRepository{
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
	e.GET("/apps/add", h.CreateApps, middlewares.Authenticated)
	e.POST("/apps/add", h.AddApps, middlewares.Authenticated)
	e.GET("/reviews", h.ListReviews, middlewares.Authenticated)
	e.GET("/keywords", h.ListKeywords, middlewares.Authenticated)
	e.GET("/keywords/add", h.CreateKeywords, middlewares.Authenticated)
	e.GET("/keyword-groups/add", h.CreateKeywordGroups, middlewares.Authenticated)
	e.POST("/keywords/add", h.AddKeywords, middlewares.Authenticated)
	e.GET("/keyword-groups", h.ListKeywordGroups, middlewares.Authenticated)
	e.GET("/analytics", h.LoadAnalyticsPage, middlewares.Authenticated)
	e.POST("/keyword-groups/add", h.AddKeywordGroups, middlewares.Authenticated)
	// AJAX listing Routes
	e.GET("/ajax/reviews/list", h.RetrieveReviews, middlewares.Authenticated)
	e.GET("/ajax/keywords/list", h.RetrieveKeywords, middlewares.Authenticated)
	e.GET("/ajax/keyword-groups/list", h.RetrieveKeywordGroups, middlewares.Authenticated)
	e.POST("/ajax/keywords/status", h.UpdateKeywordsStatus, middlewares.Authenticated)
	e.POST("/ajax/keyword-groups/status", h.UpdateKeywordGroupsStatus, middlewares.Authenticated)
	e.POST("/ajax/apps/status", h.UpdateAppsStatus, middlewares.Authenticated)
	e.GET("ajax/users/list", h.ListUsers, middlewares.Authenticated)


	e.POST("/ajax/keyword-groups/change-subscription", h.ChangeSubscriptionToKeywordGroup, middlewares.Authenticated)
	e.GET("/ajax/reviews/fetch", h.FetchReview, middlewares.Authenticated)
	e.GET("/ajax/dashboard/counts", h.GetDashboardCount, middlewares.Authenticated)
	e.GET("/ajax/analytics/counts", h.LoadAnalyticsCount, middlewares.Authenticated)

	// Dev Test Routes
	e.GET("/dev-test/verify-mongodb-queries", h.VerifyMongoDbQueries, middlewares.Authenticated)
	e.GET("/dev-test/slack", h.SendSlackMessage, middlewares.Authenticated)
	e.GET("/dev-test/review", h.FetchReview, middlewares.Authenticated)

	//User routes
	e.GET("/user", h.UsersList, middlewares.Authenticated)
	e.GET("/user/add", h.CreateUsers, middlewares.Authenticated)
	e.POST("/user/add", h.AddUser, middlewares.Authenticated)
	e.POST("/user/change-password", h.ChangePassword, middlewares.Authenticated)
	e.POST("/user/edit", h.UpdateUser, middlewares.Authenticated)

	// File server
	e.Static("/static", "public/static")

	fmt.Println("Registering Routes : DONE")

	return h
}
