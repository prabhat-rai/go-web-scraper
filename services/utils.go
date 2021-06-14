package services

import (
	"echoApp/model"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type (
	DataTableFilters struct {
		Limit int64
		Offset int64
		Search string
		SortOrder int64
		SortColumnName string
	}
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

func QueryToDataTables(c echo.Context) (dataTableFilters *DataTableFilters) {
	var sortOrder int64 = 0
	columnNumber := c.QueryParam("order[0][column]")
	startFrom, _ := strconv.ParseInt(c.QueryParam("start"), 10, 64)
	recordCount, _ := strconv.ParseInt(c.QueryParam("length"), 10, 64)

	order := c.QueryParam("order[0][dir]")
	if order == "asc" {
		sortOrder = 1
	} else {
		sortOrder = -1
	}

	return &DataTableFilters{
		Search: c.QueryParam("search[value]"),
		Offset: startFrom,
		Limit: recordCount,
		SortColumnName: c.QueryParam("columns["+columnNumber+"][name]"),
		SortOrder:  sortOrder,
	}
}

func RemoveDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func InArray(key string, haystack []string) bool {
	for _, value := range haystack {
		if value == key {
			return true
		}
	}
	return false
}
