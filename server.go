package main

import (
	"context"
	"echoApp/conf"
	"echoApp/handler"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get instance of echo framework with template setup
	e := setupFramework()

	// Database connection
	client, dbContext := connectToMongo()
	defer func() {
		clientError := client.Disconnect(dbContext)
		if clientError != nil {
			log.Fatal("DB Client error")
		}
	}()

	// Routes
	handler := registerRoutes(e, client)
	handler.Config = conf.New(client)
	scheduleFetchReviews(handler)
	// Start server
	appPort := os.Getenv("APP_PORT")
	e.Logger.Fatal(e.Start(":" + appPort))
}

func setupFramework() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())

	templateRenderer := &Template{
		templateCache: GetTemplateCache(),
	}

	e.Renderer = templateRenderer
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	fmt.Println("Framework & Template Setup : DONE")

	return e
}

func connectToMongo() (*mongo.Client, context.Context) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// DB connection via MongoDB Go Driver
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + dbHost + ":"+ dbPort +"/"))
	if err != nil {
		log.Fatal(err)
	}

	dbContext, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(dbContext)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(dbContext, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to MongoDB : DONE")
	return client, dbContext
}

func scheduleFetchReviews(h *handler.Handler) {
	schedulerConfigs := h.Config.AllSchedulerConfigs.SchedulerConfigs
	for _, schedulerConfig := range schedulerConfigs {
		//fmt.Println(schedulerConfig.Concept,schedulerConfig.Cronexpression)
		ioscron := gocron.NewScheduler(time.UTC)
		ioscron.Cron(schedulerConfig.Cronexpression).Do(h.FetchAndSaveReviews,"all",schedulerConfig.Concept)
		ioscron.StartAsync()
	}
}
