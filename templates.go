package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"time"
)

type Template struct {
	templateCache map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templateCache[name]
	if !ok {
		err := errors.New("Template not found - " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, name, data)
}


func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"noEscape": noescape,
}

// public/views/pages/*.tmpl
func GetTemplateCache() map[string]*template.Template {
	templateCache := make(map[string]*template.Template)
	templateCache["index.tmpl"] = template.Must(template.ParseFiles("public/views/pages/index.tmpl", "public/views/layouts/app.tmpl", "public/views/modules/header.tmpl", "public/views/modules/sidebar.tmpl", "public/views/modules/footer.tmpl"))
	templateCache["list.tmpl"] = template.Must(template.ParseFiles("public/views/pages/list.tmpl", "public/views/layouts/app.tmpl", "public/views/modules/header.tmpl", "public/views/modules/sidebar.tmpl", "public/views/modules/footer.tmpl"))
	templateCache["login.tmpl"] = template.Must(template.ParseFiles("public/views/pages/login.tmpl", "public/views/layouts/auth.tmpl"))
	templateCache["register.tmpl"] = template.Must(template.ParseFiles("public/views/pages/register.tmpl", "public/views/layouts/auth.tmpl"))
	templateCache["reviews.tmpl"] = template.Must(template.ParseFiles("public/views/pages/reviews.tmpl", "public/views/layouts/auth.tmpl"))

	return templateCache
}
