package services

import (
	"echoApp/model"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleDbError(err interface{}) (error error) {
	fmt.Println(err)
	return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Something is not right!!"}
}

func GetFlashMessage(c echo.Context) interface{} {
	sess, _ := session.Get("session", c)
	flashMessage := sess.Values["flash"]
	delete(sess.Values, "flash")
	sess.Save(c.Request(),c.Response())

	return flashMessage
}

func IsAuthenticated(c echo.Context) bool {
	sess, _ := session.Get("session", c)
	return sess.Values["authenticated"] == true
}

func SetFlashMessage(c echo.Context, message string) {
	SetSessionValue(c, "flash", message)
}

func SetSessionValue(c echo.Context, key string, value interface{}) {
	sess, _ := session.Get("session", c)
	sess.Values[key] = value
	sess.Save(c.Request(),c.Response())
}

func GetAuthenticatedUser(c echo.Context) *model.User {
	userName := fmt.Sprintf("%v", GetSessionValue(c, "userName"))
	userEmail := fmt.Sprintf("%v", GetSessionValue(c, "userEmail"))

	return &model.User{
		Name: userName,
		Email: userEmail,
	}
}

func GetSessionValue(c echo.Context, key string) interface{} {
	sess, _ := session.Get("session", c)
	return sess.Values[key]
}
