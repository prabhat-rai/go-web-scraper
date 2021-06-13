package handler

import (
	"echoApp/services"
	"github.com/dav009/flash"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func (h *Handler) FetchReview(c echo.Context) (err error) {
	platform := strings.ToLower(c.QueryParam("platform"))
	words := flash.NewKeywords()

	if platform == "" {
		platform = "ios"
	}

	dtf := &services.DataTableFilters{}
	keywords := h.KeywordRepository.RetrieveKeywords(dtf)

	for _, elem := range keywords.Data {
		words.Add(elem.Name)
	}

	reviews := services.FetchReview(platform, h.Config, words)
	err = h.AppReviewRepository.AddBulkReviews(reviews)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "All Ok : Fetched reviews for " + platform + ".")
}


func (h *Handler) RetrieveReviews(c echo.Context) (err error) {
	var filters = make(map[string]string)
	keywords := []string{}

	concept := c.QueryParam("concept")
	platform := c.QueryParam("platform")
	rating := c.QueryParam("rating")
	keywordGroup := c.QueryParam("keyword_groups")
	dataTableFilters := services.QueryToDataTables(c)


	if concept != "" {
		filters["concept"] = concept
	}

	if platform != "" {
		filters["platform"] = platform
	}

	if rating != "" {
		filters["rating"] = rating
	}

	if keywordGroup != "" {
		keywords = h.KeywordGroupRepository.GetKeywordsForGroup(keywordGroup)
	}

	appReviews := h.AppReviewRepository.RetrieveBulkReviews(dataTableFilters, filters, keywords)
	return c.JSON(http.StatusOK, appReviews)
}

func (h *Handler) ListReviews(c echo.Context) (err error) {
	userData := services.GetAuthenticatedUser(c)

	return c.Render(http.StatusOK, "reviews.tmpl", map[string]interface{}{
		"name": userData.Name,
		"reviews": nil,
		"concepts": h.Config.AllApps.Apps,
		"platforms" : []string{"ios", "android"},
		"ratings" : []int{1,2,3,4,5},
		"keyword_groups" : []string{"SHUKRAN"},
	})
}