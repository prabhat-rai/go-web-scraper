package handler

import (
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) ListKeywordGroups(c echo.Context) (err error) {
	userData := services.GetAuthenticatedUser(c)

	return c.Render(http.StatusOK, "keyword_groups.tmpl", map[string]interface{}{
		"name": userData.Name,
	})
}

func (h *Handler) RetrieveKeywordGroups(c echo.Context) (err error) {
	dataTableFilters := services.QueryToDataTables(c)
	activeRecords := c.QueryParam("active")

	userData := services.GetAuthenticatedUser(c)
	keywordGroups := h.KeywordGroupRepository.RetrieveKeywordGroups(dataTableFilters, userData, activeRecords)
	return c.JSON(http.StatusOK, keywordGroups)
}

func (h *Handler) ChangeSubscriptionToKeywordGroup (c echo.Context) (err error) {
	keyGroupId := c.FormValue("id")
	subscriptionStatus := c.FormValue("subscription")
	userData := services.GetAuthenticatedUser(c)

	result := h.KeywordGroupRepository.UpdateSubscriptionForUser(keyGroupId, subscriptionStatus, userData.Email)
	return c.JSON(http.StatusOK, result)
}