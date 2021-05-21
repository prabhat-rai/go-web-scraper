package main

import (
	"echoApp/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/mgo.v2"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.tmpl")),
	}
	e.Renderer = t

	// Database connection
	db, err := mgo.Dial("localhost")
	if err != nil {
		e.Logger.Fatal(err)
	}

	dbName := "lmg_reviews"

	// Create indices
	if err = db.Copy().DB(dbName).C("app_reviews").EnsureIndex(mgo.Index{
		Key:    []string{"review_id"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	// Create indices on users
	if err = db.Copy().DB(dbName).C("users").EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	// Initialize handler
	h := &handler.Handler{DB: db, DbName: dbName}

	// Routes
	e.GET("/", h.Home)
	e.GET("/dev-test/verify-mongodb-queries", h.VerifyMongoDbQueries)
	e.GET("/dev-test/review", h.FetchReview)

	// Start server
	e.Logger.Fatal(e.Start(":1200"))
}
