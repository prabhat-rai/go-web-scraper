package handler

import (
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"fmt"
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
	appReviews := h.AppReviewRepository.RetrieveBulkReviews()
	fmt.Printf("%+v", appReviews)
	return c.Render(http.StatusOK, "reviews.tmpl", map[string]interface{}{
		"reviews": appReviews.AppReview,
	})
}