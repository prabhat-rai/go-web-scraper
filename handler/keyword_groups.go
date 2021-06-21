package handler

import (
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"echoApp/model"
	"net/http"
	"log"
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

func (h *Handler) AddKeywordGroups(c echo.Context) (err error) {
	log.Println("hello")

	active := true
	if c.FormValue("active") == "true" {
		active = true
	} else{
		active = false
	}
	log.Println(c.FormValue("keywords[]"))

	keyword_group := &model.KeywordGroup{
		GroupName: c.FormValue("keyword_group"),
		Active:active,
	}

	err = c.Bind(keyword_group)
	if err != nil {
		//place holder to render register page with error message
		return c.Render(http.StatusOK, "create_keywords.tmpl", map[string]interface{}{
			"Flash": "Something went wrong!! Please Try Again.",
		})
	}

	err = h.KeywordGroupRepository.CreateKeywordGroup(keyword_group)
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