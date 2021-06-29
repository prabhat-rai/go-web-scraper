package handler

import (
	"echoApp/model"
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) Home(c echo.Context) (err error) {
	conceptCountsLastWeek := h.AppReviewRepository.CountReviews("concept", 7, "days")
	platformCountsLastWeek := h.AppReviewRepository.CountReviews("platform", 7, "days")
	conceptCountsLastMonth := h.AppReviewRepository.CountReviews("concept", 1, "months")
	platformCountsLastMonth := h.AppReviewRepository.CountReviews("platform", 1, "months")

	commonData := services.GetCommonDataForTemplates(c)

	return c.Render(http.StatusOK, "index.tmpl", map[string]interface{}{
		"commonData" : commonData,
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
		services.SetFlashMessage(c, "Username or Password is incorrect.")
		return c.Redirect(http.StatusFound, "/login")
	}

	user,err :=h.UserRepository.AuthenticateUser(u.Email, u.Password)
	if err != nil {
		services.SetFlashMessage(c, "Username or Password is incorrect.")
		return c.Redirect(http.StatusFound, "/login")
	}

	services.SetSessionValue(c, "authenticated", true)
	services.SetSessionValue(c, "userName", user.Name)
	services.SetSessionValue(c, "userEmail", user.Email)
	services.SetSessionValue(c, "role", user.Role)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) Logout(c echo.Context) (err error) {
	services.SetSessionValue(c, "authenticated", false)
	services.SetSessionValue(c, "userName", false)
	services.SetSessionValue(c, "userEmail", false)
	services.SetSessionValue(c, "role", false)
	return c.Redirect(http.StatusSeeOther, "/login")
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
		services.SetFlashMessage(c, "Something went wrong!! Please Try Again.")
		return c.Redirect(http.StatusFound, "/register")
	}

	err = h.UserRepository.CreateUser(user)
	if err != nil {
		services.SetFlashMessage(c, "Something went wrong!! Please Try Again.")
		return c.Redirect(http.StatusFound, "/register")
	}

	services.SetSessionValue(c, "authenticated", true)
	services.SetSessionValue(c, "userName", user.Name)
	services.SetSessionValue(c, "userEmail", user.Email)
	services.SetSessionValue(c, "role", user.Role)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) ListUsers(c echo.Context) (err error) {
	dataTableFilters := services.QueryToDataTables(c)
	searchWord := c.QueryParam("searchWord")
	users, err := h.UserRepository.ListUsers(dataTableFilters,searchWord)
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) UsersList(c echo.Context) (err error) {
	commonData := services.GetCommonDataForTemplates(c)

	return c.Render(http.StatusOK, "users_list.tmpl", map[string]interface{}{
		"commonData" : commonData,
	})
}

func (h *Handler) CreateUsers(c echo.Context) (err error) {
	commonData := services.GetCommonDataForTemplates(c)
	return c.Render(http.StatusOK, "create_users.tmpl", map[string]interface{}{
		"commonData" : commonData,
	})
}

func (h *Handler) AddUser(c echo.Context) (err error) {
	user := &model.User{
		Name: c.FormValue("name"),
		Email: c.FormValue("email"),
		Phone: c.FormValue("phone"),
		Role: c.FormValue("role"),
	}
	err = c.Bind(user)
	if err != nil {
		services.SetFlashMessage(c, "Something went wrong!! Please Try Again.")
		return err
	}
	//Set default Password
	user.Password = h.Config.ConfigProps.DefaultPassword
	err = h.UserRepository.CreateUser(user)
	if err != nil {
		log.Println(err)
		return err
	}
	services.SetSuccessMessage(c, "User added successfully!")
	return c.Redirect(http.StatusFound, "/user")
}

func (h *Handler) ChangePassword(c echo.Context) (err error) {
	authUser := services.GetAuthenticatedUser(c)
	password := c.FormValue("password")

	err = h.UserRepository.ChangePassword(authUser.Email,password)
	if err != nil {
		log.Println("Error while updating the password", err)
		return err
	}
	return err
}

func (h *Handler) UpdateUser(c echo.Context) (err error) {
	id, err := primitive.ObjectIDFromHex(c.FormValue("id"))
	name := c.FormValue("name")
	email := c.FormValue("email")
	role := c.FormValue("role")
	phone := c.FormValue("phone")

	if err != nil {
		log.Println(err)
		return err
	}

	user := &model.User{
		ID: id,
		Name: name,
		Email: email,
		Role: role,
		Phone: phone,
	}

	err = h.UserRepository.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Something went wrong!! Please Try Again.")
	}

	services.SetSuccessMessage(c, "User updated successfully!")
	return c.Redirect(http.StatusFound, "/user")
}