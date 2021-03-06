package services

import (
	"echoApp/model"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	CommonDataForTemplates struct {
		Path 		string
		UserName 	string
		Role		string
		User 		*model.User
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

func SetSuccessMessage(c echo.Context, message string) {
	SetSessionValue(c, "message", message)
}

func SetSessionValue(c echo.Context, key string, value interface{}) {
	sess, _ := session.Get("session", c)
	sess.Values[key] = value
	sess.Save(c.Request(),c.Response())
}

func GetAuthenticatedUser(c echo.Context) *model.User {
	userName := fmt.Sprintf("%v", GetSessionValue(c, "userName"))
	userEmail := fmt.Sprintf("%v", GetSessionValue(c, "userEmail"))
	role := fmt.Sprintf("%v", GetSessionValue(c, "role"))
	return &model.User{
		Name: userName,
		Email: userEmail,
		Role: role,
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

func GetKeyBasedCount(aggregateResultFromDb []bson.M) map[string]int {
	keyBasedCount := make(map[string]int)
	for _, data := range aggregateResultFromDb {
		aggregator := fmt.Sprintf("%s", data["_id"])
		count := fmt.Sprintf("%d", data["count"])
		keyBasedCount[aggregator], _ = strconv.Atoi(count)
	}

	return keyBasedCount
}

func GetKeyBasedCountForDailyBasis(aggregateResultFromDb []bson.M, aggregator string) map[string]map[string]int {
	keyBasedCount := make(map[string]map[string]int)

	for _, data := range aggregateResultFromDb {
		aggregator := fmt.Sprintf("%s", data["_id"].(primitive.M)[aggregator])
		dateKey := fmt.Sprintf("%s", data["_id"].(primitive.M)["review_date"])
		count, _ := strconv.Atoi(fmt.Sprintf("%d", data["count"]))

		if keyBasedCount[dateKey] == nil {
			keyBasedCount[dateKey] = make(map[string]int)
		}

		keyBasedCount[dateKey][aggregator] = count
	}

	return keyBasedCount
}

func GetCommonDataForTemplates(c echo.Context) CommonDataForTemplates {

	path := c.Path()
	platform := c.QueryParam("platform")

	if platform != "" {
		path += "?platform=" + platform
	}

	authUser := GetAuthenticatedUser(c)
	return CommonDataForTemplates{
		UserName: authUser.Name,
		Role: authUser.Role,
		Path: path,
	}
}
