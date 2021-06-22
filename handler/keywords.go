package handler

import (
	"echoApp/services"
	"echoApp/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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

func (h *Handler) CreateKeywords(c echo.Context) (err error) {
	userData := services.GetAuthenticatedUser(c)
	return c.Render(http.StatusOK, "create_keywords.tmpl", map[string]interface{}{
		"name": userData.Name,		
	})
}

func (h *Handler) AddKeywords(c echo.Context) (err error) {
	active := false
	if c.FormValue("active") == "true" {
		active = true
	}
	keyword := &model.Keyword{
		Name: c.FormValue("keyword_name"),
		Active:active,
	}

	err = h.KeywordRepository.CreateKeyword(keyword)
	if err != nil {
		return c.Render(http.StatusOK, "create_keywords.tmpl", map[string]interface{}{
			"Flash": "Something went wrong!! Please Try Again.",
		})
	}
	userData := services.GetAuthenticatedUser(c)
	return c.Render(http.StatusOK, "keywords.tmpl", map[string]interface{}{
		"message": "Keyword added successfully",
		"name": userData.Name,
	})

	return nil
}

func (h *Handler) UpdateKeywordsStatus(c echo.Context) (err error) {
	active := false
	id, err := primitive.ObjectIDFromHex(c.QueryParam("id"))
	if err != nil {
	log.Fatal(err)
	}
	if c.QueryParam("active") == "true" {
		active = true
	}

	keyword := &model.Keyword{
		ID: id,
		Active: active,
	}

	err = h.KeywordRepository.UpdateActiveStatus(keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Something went wrong!! Please Try Again.")
	}

	return c.JSON(http.StatusOK, "Updated")
}




