package services

import (
	"echoApp/model"
	"encoding/xml"
	"github.com/dav009/flash"
	"github.com/n0madic/google-play-scraper/pkg/reviews"
	"github.com/n0madic/google-play-scraper/pkg/store"
	"log"
	"net/http"
	"strconv"
)

type Feed struct {
	Id string 			`xml:"id"`
	Title string 		`xml:"title"`
	Updated string 		`xml:"updated"`
	Entry []Entry 		`xml:"entry"`
}

type Entry struct {
	Updated string 		`xml:"updated"`
	Id string 			`xml:"id"`
	Title string 		`xml:"title"`
	Content []Content 	`xml:"content"`
	Rating string 		`xml:"rating"`
	Version string 		`xml:"version"`
	Author Author 		`xml:"author"`
}

type Content struct {
	Type    string 		`xml:"type,attr"`
	Data	string 		`xml:",chardata"`
}

type Author struct {
	Name string 		`xml:"name"`
	Uri string 			`xml:"uri"`
}

func LoadAndroidReviews(id string, appName string, keywords flash.Keywords, latestReviewId string) []*model.AppReview {
	r := reviews.New("com.landmarkgroup." + id, reviews.Options{
		Number: 50,
		Sorting: store.SortNewest,
		Language: "en",
	})

	err := r.Run()
	if err != nil {
		panic(err)
	}

	var appReviews []*model.AppReview
	for _, review := range r.Results {
		foundKeywords := RemoveDuplicateValues(keywords.Extract(review.Text))

		if latestReviewId == review.ID {
			break
		}

		appReviews = append(appReviews, &model.AppReview{
			ReviewId: review.ID,
			ReviewDate: review.Timestamp.String(),
			UserName: review.Reviewer,
			Title: "",
			Description: review.Text,
			Rating: strconv.Itoa(review.Score),
			CreatedAt: review.Timestamp.String(),
			UpdatedAt: review.Timestamp.String(),
			Platform: "android",
			Version: review.Version,
			Concept: appName,
			Keywords: foundKeywords,
		})
	}

	return appReviews
}

func LoadIosReviews(id string, appName string, keywords flash.Keywords, latestReviewId string) []*model.AppReview {

	urlPrefix := "https://itunes.apple.com/ae/rss/customerreviews/id="
	urlSuffix := "/page=1/sortBy=mostRecent/xml"
	url := urlPrefix + id + urlSuffix

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error GET: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	feed := Feed{}

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&feed)
	if err != nil {
		log.Printf("Error Decode: %v\n", err)
		return nil
	}

	var appReviews []*model.AppReview
	for _, entry := range feed.Entry {
		foundKeywords := RemoveDuplicateValues(keywords.Extract(entry.Content[0].Data))

		if latestReviewId == entry.Id {
			break
		}

		appReviews = append(appReviews, &model.AppReview{
			ReviewId: entry.Id,
			ReviewDate: entry.Updated,
			UserName: entry.Author.Name,
			Title: entry.Title,
			Description: entry.Content[0].Data,
			Rating: entry.Rating,
			CreatedAt: entry.Updated,
			UpdatedAt: entry.Updated,
			Platform: "ios",
			Version: entry.Version,
			Concept: appName,
			Keywords: foundKeywords,
		})
	}

	return appReviews
}