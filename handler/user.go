package handler

import (
	"echoApp/model"
	"echoApp/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Home(c echo.Context) (err error) {
	userData := services.GetAuthenticatedUser(c)

	return c.Render(http.StatusOK, "index.tmpl", map[string]interface{}{
		"name": userData.Name,
	})
}

func (h *Handler) List(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "list.tmpl", map[string]interface{}{
		"name": "Admin",
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

	fmt.Printf("User authenticated %s",user.Email)
	services.SetSessionValue(c, "authenticated", true)
	services.SetSessionValue(c, "userName", user.Name)
	services.SetSessionValue(c, "userEmail", user.Email)
	c.Redirect(http.StatusSeeOther, "/")
	return nil
}

func (h *Handler) Logout(c echo.Context) (err error) {
	services.SetSessionValue(c, "authenticated", false)
	services.SetSessionValue(c, "userName", false)
	services.SetSessionValue(c, "userEmail", false)
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

	fmt.Printf("%+v\n", user)
	//fmt.Printf("%v", c.FormValue("name"))
	err = c.Bind(user)
	if err != nil {
		//place holder to render register page with error message
		return c.Render(http.StatusOK, "register.tmpl", map[string]interface{}{
			"Flash": "Something went wrong!! Please Try Again.",
		})
	}

	fmt.Printf("%+v\n", user)
	err = h.UserRepository.CreateUser(user)
	if err != nil {
		return c.Render(http.StatusOK, "register.tmpl", map[string]interface{}{
			"Flash": "Something went wrong!! Please Try Again.",
		})
	}

	services.SetSessionValue(c, "authenticated", true)
	services.SetSessionValue(c, "userName", user.Name)
	services.SetSessionValue(c, "userEmail", user.Email)
	c.Redirect(http.StatusSeeOther, "/")

	return nil
}