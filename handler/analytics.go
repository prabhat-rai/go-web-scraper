package handler

import (
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *Handler) LoadAnalyticsPage(c echo.Context) (err error) {
	userData := services.GetAuthenticatedUser(c)

	return c.Render(http.StatusOK, "analytics.tmpl", map[string]interface{}{
		"name": userData.Name,
	})
}

func (h *Handler) LoadAnalyticsCount(c echo.Context) (err error) {
	noOfDays, err := strconv.Atoi(c.QueryParam("days"))

	if err != nil {
		noOfDays = 7
	}

	conceptWiseCount := h.AppReviewRepository.DateWiseReviews("concept", noOfDays)
	conceptWiseCountNew := services.GetKeyBasedCountForDailyBasis(conceptWiseCount, "concept")

	platformWiseCount := h.AppReviewRepository.DateWiseReviews("platform", noOfDays)
	platformWiseCountNew := services.GetKeyBasedCountForDailyBasis(platformWiseCount, "platform")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"concept" : conceptWiseCountNew,
		"platform" : platformWiseCountNew,
		"noOfDays" : noOfDays + 1,
	})
}

func (h *Handler) GetDashboardCount(c echo.Context) (err error) {
	platformCountsLastWeek := h.AppReviewRepository.CountReviews("platform", 7, "days")
	platformCountsLastMonth := h.AppReviewRepository.CountReviews("platform", 1, "months")
	conceptCountsLastWeek := h.AppReviewRepository.CountReviews("concept", 7, "days")
	conceptCountsLastMonth := h.AppReviewRepository.CountReviews("concept", 1, "months")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"concept" : map[string]interface{}{
			"week" : services.GetKeyBasedCount(conceptCountsLastWeek),
			"week_orig" : conceptCountsLastWeek,
			"month" : services.GetKeyBasedCount(conceptCountsLastMonth),
		},
		"platform" : map[string]interface{}{
			"week" : services.GetKeyBasedCount(platformCountsLastWeek),
			"month" : services.GetKeyBasedCount(platformCountsLastMonth),
		},
	})
}