package handler

import (
	"echoApp/conf"
	"echoApp/model"
	"echoApp/services"
	"github.com/dav009/flash"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) FetchReview(c echo.Context) (err error) {
	platform := strings.ToLower(c.QueryParam("platform"))
	concept := strings.ToLower(c.QueryParam("concept"))



	if platform == "" {
		platform = "all"
	}

	h.FetchAndSaveReviews( platform, concept)

	return c.JSON(http.StatusOK, "All Ok : Fetched reviews for "+platform+" platform.")
}

func (h *Handler) FetchAndSaveReviews( platform string, concept string) {
	words := flash.NewKeywords()
	dtf := &services.DataTableFilters{}
	keywords := h.KeywordRepository.RetrieveKeywords(dtf)

	for _, elem := range keywords.Data {
		words.Add(elem.Name)
	}

	if platform == "ios" || platform == "all" {
		go h.fetchReviewForApp(concept, "ios", h.Config.AllApps, words)
	}

	if platform == "android" || platform == "all" {
		go h.fetchReviewForApp(concept, "android", h.Config.AllApps, words)
	}
}

func (h *Handler) fetchReviewForApp(concept string, platform string, config conf.AllApps, words flash.Keywords) bool {
	var reviews []*model.AppReview

	for _, elem := range config.Apps {
		if concept == "" || concept == elem.Name {
			latestReviewId := h.AppReviewRepository.GetLatestReviewId(platform, elem.Name)
			log.Printf("STARTING : %s Reviews for %s \n\n", platform, elem.Name)

			if platform == "ios" {
				reviews = append(reviews, services.LoadIosReviews(elem.IosAppId, elem.Name, words, latestReviewId)...)
			} else {
				reviews = append(reviews, services.LoadAndroidReviews(elem.GoogleAppId, elem.Name, words, latestReviewId)...)
			}
			log.Printf("DONE : %s Reviews for %s \n\n", platform, elem.Name)
		}
	}

	if len(reviews) > 0 {
		insertedIDs, err := h.AppReviewRepository.AddBulkReviews(reviews)

		if err != nil {
			log.Println(err)
			return false
		}

		h.TriggerMailToSubscribers(insertedIDs)
	}


	return true
}

func (h *Handler) TriggerMailToSubscribers(insertedIDs interface{}) {
	//log.Println("Inserted Docs: ", insertedIDs)


	// Get All KeywordGroups where subscriber is not empty
	keywordGroups := h.KeywordGroupRepository.GetGroupsWithActiveSubscribers()

	// Loop through the fetched groups and run query with Keyword present in group + inserted Id in the reviews collection
	for _, elem := range keywordGroups {
		//elem.Keywords

		reviews := h.AppReviewRepository.GetReviewsWithMatchingKeywords(elem.Keywords, insertedIDs)

		// If we have reviews send a mail to subscribers.
		if len(reviews) > 0 {
			keywordsString := strings.Join(elem.Keywords, ",")
			mailConfig := h.Config.ConfigProps.MailConfig
			go services.SendMailForNewReviews(elem.Subscribers, reviews, elem.GroupName, keywordsString, mailConfig)
		}
	}
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