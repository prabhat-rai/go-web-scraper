package handler

import (
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"echoApp/model"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
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
	active := true
	if c.FormValue("active") == "true" {
		active = true
	} else{
		active = false
	}
	keywords := strings.Split(c.FormValue("keywords"), ",")
	keyword_group := &model.KeywordGroup{
		GroupName: c.FormValue("keyword_group"),
		Active:active,
		Keywords:keywords,
	}

	err = h.KeywordGroupRepository.CreateKeywordGroup(keyword_group)
	if err != nil {
		return c.Render(http.StatusOK, "register.tmpl", map[string]interface{}{
			"Flash": "Something went wrong!! Please Try Again.",
		})
	}
	 c.Redirect(http.StatusSeeOther, "/keywords/add")

	return nil
}

func (h *Handler) UpdateKeywordGroupsStatus(c echo.Context) (err error) {
	active := false
	id, err := primitive.ObjectIDFromHex(c.QueryParam("id"))
	if err != nil {
	panic(err)
	}
	if c.QueryParam("active") == "true" {
		active = true
	}

	keyword_group := &model.KeywordGroup{
		ID: id,
		Active: !active,
	}

	err = h.KeywordGroupRepository.UpdateKeywordGroupStatus(keyword_group)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Something went wrong!! Please Try Again.")
	}

	return c.JSON(http.StatusOK, "Updated")
}