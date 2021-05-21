package services

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2"
	"net/http"
)

func HandleDbError(err interface{}) (error error) {
	if err == mgo.ErrNotFound {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Record not found"}
	}

	fmt.Println(err)
	return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Something is not right!!"}
}
