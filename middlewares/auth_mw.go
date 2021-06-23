package middlewares

import (
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Proceed if user is logged in
		if services.IsAuthenticated(c) {
			return next(c)
		}

		// Redirect to login page if session is not active
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

func Guest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// If user is not logged in then load guest page
		if !services.IsAuthenticated(c) {
			return next(c)
		}

		// If user is logged in load the homepage
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}