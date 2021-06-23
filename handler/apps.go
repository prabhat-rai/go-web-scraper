package handler

import (
	"echoApp/conf"
	"echoApp/model"
	"echoApp/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) AppsList(c echo.Context) (err error) {
	commonData := services.GetCommonDataForTemplates(c)
	appDetails := conf.GetAppsConfig(h.DB, false)
	return c.Render(http.StatusOK, "apps_list.tmpl", map[string]interface{}{
		"commonData" : commonData,
		"apps": appDetails.Apps,
	})
}

func (h *Handler) CreateApps(c echo.Context) (err error) {
	commonData := services.GetCommonDataForTemplates(c)
	return c.Render(http.StatusOK, "create_apps.tmpl", map[string]interface{}{
		"commonData" : commonData,
	})
}
func (h *Handler) AddApps(c echo.Context) (err error) {
	active := false
	if c.FormValue("active") == "true" {
		active = true
	}
	app := &model.Apps{
		Name: c.FormValue("app_name"),
		GoogleAppId: c.FormValue("google_app_id"),
		IosAppId: c.FormValue("ios_app_id"),	
		Active: active,
	}

	err = h.AppRepository.CreateApp(app)
	if err != nil {
		services.SetFlashMessage(c, "Something went wrong!! Please Try Again.")
		return c.Redirect(http.StatusFound, "/apps/add")
	}

	services.SetSuccessMessage(c, "App successfully created!")
	return c.Redirect(http.StatusFound, "/apps")
}

func (h *Handler) UpdateAppsStatus(c echo.Context) (err error) {
	id, err := primitive.ObjectIDFromHex(c.QueryParam("id"))
	active := false
	if c.QueryParam("active") == "true" {
		active = true
	}

	app := &model.Apps{
		ID: id,
		Active: active,
	}

	err = h.AppRepository.UpdateActiveStatus(app)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Something went wrong!! Please Try Again.")
	}

	return c.JSON(http.StatusOK, "Updated")
}