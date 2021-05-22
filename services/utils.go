package services

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleDbError(err interface{}) (error error) {
	fmt.Println(err)
	return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Something is not right!!"}
}
