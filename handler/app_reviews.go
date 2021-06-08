package handler

import (
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func (h *Handler) FetchReview(c echo.Context) (err error) {
	platform := strings.ToLower(c.QueryParam("platform"))

	if platform == "" {
		platform = "ios"
	}

	reviews := services.FetchReview(platform, h.Config)
	err = h.AppReviewRepository.AddBulkReviews(reviews)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "All Ok : Fetched reviews for " + platform + ".")
}

func (h *Handler) RetrieveReviews(c echo.Context) (err error) {

	
	err = h.AppReviewRepository.RetrieveBulkReviews()

	if err != nil {
		return err
	}
	return nil
	//return c.JSON(http.StatusOK)
}
