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

	keywordGroups := h.KeywordGroupRepository.RetrieveKeywordGroups(dataTableFilters)
	return c.JSON(http.StatusOK, keywordGroups)
}
