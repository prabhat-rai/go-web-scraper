package handler

import (
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) ListKeywords(c echo.Context) (err error) {
	userData := services.GetAuthenticatedUser(c)

	return c.Render(http.StatusOK, "keywords.tmpl", map[string]interface{}{
		"name": userData.Name,
	})
}

func (h *Handler) RetrieveKeywords(c echo.Context) (err error) {
	dataTableFilters := services.QueryToDataTables(c)

	keywords := h.KeywordRepository.RetrieveKeywords(dataTableFilters)
	return c.JSON(http.StatusOK, keywords)
}
