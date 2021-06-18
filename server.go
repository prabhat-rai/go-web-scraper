package main

import (
	"context"
	"echoApp/conf"
	"echoApp/handler"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {
	configProps, err := conf.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Get instance of echo framework with template setup
	e := setupFramework()

	// Database connection
	client, dbContext := connectToMongo(configProps)
	defer func() {
		clientError := client.Disconnect(dbContext)
		if clientError != nil {
			log.Fatal("DB Client error")
		}
	}()

	// Routes
	handler := registerRoutes(e, client,configProps.DB_DATABASE)
	handler.Config = conf.New(client,configProps)
	scheduleFetchReviews(handler,configProps.SCHEDULER_CONFIGS)
	// Start server
	//appPort := os.Getenv("APP_PORT")
	appPort := configProps.APP_PORT
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

func connectToMongo(configProps conf.ConfigProps) (*mongo.Client, context.Context) {
	dbHost := configProps.DB_HOST
	dbPort := configProps.DB_PORT

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

func scheduleFetchReviews(h *handler.Handler, schedulerConfigs map[string]string) {
	for concept, cronexpression := range schedulerConfigs {
		ioscron := gocron.NewScheduler(time.UTC)
		ioscron.Cron(cronexpression).Do(h.FetchAndSaveReviews,"all",concept)
		ioscron.StartAsync()
	}
}
