package handler

import (
	"echoApp/model"
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Home(c echo.Context) (err error) {
	conceptCountsLastWeek := h.AppReviewRepository.CountReviews("concept", 7, "days")
	platformCountsLastWeek := h.AppReviewRepository.CountReviews("platform", 7, "days")
	conceptCountsLastMonth := h.AppReviewRepository.CountReviews("concept", 1, "months")
	platformCountsLastMonth := h.AppReviewRepository.CountReviews("platform", 1, "months")

	userData := services.GetAuthenticatedUser(c)

	return c.Render(http.StatusOK, "index.tmpl", map[string]interface{}{
		"name": userData.Name,
		"role" : userData.Role,
		"concept" : map[string]interface{}{
			"week" : services.GetKeyBasedCount(conceptCountsLastWeek),
			"month" : services.GetKeyBasedCount(conceptCountsLastMonth),
		},
		"platform" : map[string]interface{}{
			"week" : services.GetKeyBasedCount(platformCountsLastWeek),
			"month" : services.GetKeyBasedCount(platformCountsLastMonth),
		},
	})
}


func (h *Handler) LoginForm(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "login.tmpl", map[string]interface{}{
		"Flash": services.GetFlashMessage(c),
	})
}

func (h *Handler) Login(c echo.Context) error {
	u := &model.User{
		Email: c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	err := c.Bind(u)
	if err != nil {
		return c.Render(http.StatusOK, "login.tmpl", map[string]interface{}{
			"Flash": services.GetFlashMessage(c),
		})
	}

	user,err :=h.UserRepository.AuthenticateUser(u.Email, u.Password)
	if err != nil {
		services.SetFlashMessage(c, "Username or Password is incorrect.")
		c.Redirect(http.StatusFound, "/login")
		return nil
	}

	services.SetSessionValue(c, "authenticated", true)
	services.SetSessionValue(c, "userName", user.Name)
	services.SetSessionValue(c, "userEmail", user.Email)
	services.SetSessionValue(c, "role", user.Role)
	c.Redirect(http.StatusSeeOther, "/")
	return nil
}

func (h *Handler) Logout(c echo.Context) (err error) {
	services.SetSessionValue(c, "authenticated", false)
	services.SetSessionValue(c, "userName", false)
	services.SetSessionValue(c, "userEmail", false)
	services.SetSessionValue(c, "role", false)
	c.Redirect(http.StatusSeeOther, "/login")
	return err
}

func (h *Handler) RegisterForm(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "register.tmpl", map[string]interface{}{
		"Flash": services.GetFlashMessage(c),
	})
}

func (h *Handler) Register(c echo.Context) (err error) {
	// @TODO : VALIDATE THE REQUEST
	user := &model.User{
		Name: c.FormValue("name"),
		Email: c.FormValue("email"),
		Password: c.FormValue("password"),
		Phone: c.FormValue("phone"),
	}

	err = c.Bind(user)
	if err != nil {
		//place holder to render register page with error message
		return c.Render(http.StatusOK, "register.tmpl", map[string]interface{}{
			"Flash": "Something went wrong!! Please Try Again.",
		})
	}

	err = h.UserRepository.CreateUser(user)
	if err != nil {
		return c.Render(http.StatusOK, "register.tmpl", map[string]interface{}{
			"Flash": "Something went wrong!! Please Try Again.",
		})
	}

	services.SetSessionValue(c, "authenticated", true)
	services.SetSessionValue(c, "userName", user.Name)
	services.SetSessionValue(c, "userEmail", user.Email)
	services.SetSessionValue(c, "role", user.Role)
	c.Redirect(http.StatusSeeOther, "/")

	return nil
}


func (h *Handler) ListUsers(c echo.Context) (err error) {
	userData := services.GetAuthenticatedUser(c)
	return c.Render(http.StatusOK, "keywords.tmpl", map[string]interface{}{
		"name": userData.Name,
		"role" : userData.Role,
	})
}

func (h *Handler) AddUser(c echo.Context) (err error) {
return nil
}