package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"path/filepath"
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

func GetTemplateCache() map[string]*template.Template {
	templateCache := make(map[string]*template.Template)
	pages, _ := filepath.Glob(filepath.Join("public/views/pages/", "*.tmpl"))

	for _, page := range pages {
		name := filepath.Base(page)

		if name == "login.tmpl" || name == "register.tmpl" {
			templateCache[name] = template.Must(
				template.New(name).Funcs(functions).ParseFiles(
					page,
					"public/views/layouts/auth.tmpl",
				),
			)
		} else {
			templateCache[name] = template.Must(
				template.New(name).Funcs(functions).ParseFiles(
					page,
					"public/views/layouts/app.tmpl",
					"public/views/modules/header.tmpl",
					"public/views/modules/sidebar.tmpl",
					"public/views/modules/footer.tmpl",
				),
			)
		}
	}

	return templateCache
}
