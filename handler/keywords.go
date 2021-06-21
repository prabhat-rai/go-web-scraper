package handler

import (
	"echoApp/services"
	"echoApp/model"
	"github.com/labstack/echo/v4"
	"net/http"
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
		"page": c.QueryParam("page"),
		
	})
}

func (h *Handler) AddKeywords(c echo.Context) (err error) {
	log.Println(c.FormValue("active"))
	active := true
	if c.FormValue("active") == "true" {
		log.Println("TRUEEE")
		active = true
	} else{
		log.Println("FALSEE")
		active = false
	}
	keyword := &model.Keyword{
		Name: c.FormValue("keyword_name"),
		Active:active,
	}

	err = c.Bind(keyword)
	if err != nil {
		//place holder to render register page with error message
		return c.Render(http.StatusOK, "create_keywords.tmpl", map[string]interface{}{
			"Flash": "Something went wrong!! Please Try Again.",
		})
	}

	err = h.KeywordRepository.CreateKeyword(keyword)
	if err != nil {
		return c.Render(http.StatusOK, "register.tmpl", map[string]interface{}{
			"Flash": "Something went wrong!! Please Try Again.",
		})
	}

	// services.SetSessionValue(c, "authenticated", true)
	// services.SetSessionValue(c, "userName", user.Name)
	// services.SetSessionValue(c, "userEmail", user.Email)
	 c.Redirect(http.StatusSeeOther, "/keywords/add")

	return nil
}




