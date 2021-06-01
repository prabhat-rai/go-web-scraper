package main

import (
	"context"
	"echoApp/handler"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"io"
	"time"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
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
	registerRoutes(e, client)

	// Start server
	e.Logger.Fatal(e.Start(":1200"))
}

func setupFramework() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.tmpl")),
	}
	e.Renderer = t
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	fmt.Println("Framework & Template Setup : DONE")

	return e
}

func connectToMongo() (*mongo.Client, context.Context) {
	// DB connection via MongoDB Go Driver
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
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

func registerRoutes(e *echo.Echo, client *mongo.Client) {
	// Connect to DB
	database := client.Database("lmg_reviews")

	// Initialize handler
	h := &handler.Handler{DB: database}
	// Register Routes
	e.GET("/", h.Home)
	e.GET("/dev-test/verify-mongodb-queries", h.VerifyMongoDbQueries)
	e.GET("/dev-test/review", h.FetchReview)
	e.POST("/login", h.Login)
	e.GET("/login", h.LoginForm)
	e.GET("/register", h.RegisterForm)
	e.POST("/register", h.Register)
	e.POST("/logout", h.Logout)
	e.GET("/logout", h.Logout)

	fmt.Println("Registering Routes : DONE")
}
