package services

import (
	"echoApp/conf"
	"encoding/xml"
	"fmt"
	"github.com/n0madic/google-play-scraper/pkg/reviews"
	"github.com/n0madic/google-play-scraper/pkg/store"
	"net/http"
)

type Feed struct {
	Id string `xml:"id"`
	Title string `xml:"title"`
	Updated string `xml:"updated"`
	Entry []Entry `xml:"entry"`
}

type Entry struct {
	Updated string `xml:"updated"`
	Id string `xml:"id"`
	Title string `xml:"title"`
	Content string `xml:"content"`
	Rating string `xml:"rating"`
	Version string `xml:"version"`
	Author Author `xml:"author"`
}

type Author struct {
	Name string `xml:"name"`
	Uri string `xml:"uri"`
}

func FetchReview(platform string, conf *conf.Config) {
	if platform == "ios" {
		FetchIosReviewsForAllApps(conf.IosApps)
	} else {
		FetchAndroidReviewsForAllApps(conf.AndroidApps)
	}

}

func FetchIosReviewsForAllApps (config conf.AllIosApps) {
	for _, elem := range config.Apps {
		fmt.Printf("STARTING : IOS Reviews for %s \n\n", elem.Name)
		LoadIosReviews(elem.AppId)
		fmt.Printf("DONE : IOS Reviews for %s \n\n", elem.Name)
	}

	return
}

func FetchAndroidReviewsForAllApps(config conf.AllAndroidApps) {
	for _, elem := range config.Apps {
		fmt.Printf("STARTING : ANDROID Reviews for %s \n\n", elem.Name)
		LoadAndroidReviews(elem.GoogleAppId)
		fmt.Printf("DONE : ANDROID Reviews for %s \n\n", elem.Name)
	}

	return
}

func LoadAndroidReviews(id string) {
	r := reviews.New("com.landmarkgroup." + id, reviews.Options{
		Number: 50,
		Sorting: store.SortNewest,
		Language: "en",
	})

	err := r.Run()
	if err != nil {
		panic(err)
	}

	for i, review := range r.Results {
		fmt.Printf("%v. %v gave rating %v and wrote: %v on %v\n\n", i+1, review.Reviewer, review.Score, review.Text, review.Timestamp)
	}
}

func LoadIosReviews(id string) {
	urlPrefix := "https://itunes.apple.com/ae/rss/customerreviews/id="
	urlSuffix := "/page=1/sortBy=mostRecent/xml"
	url := urlPrefix + id + urlSuffix

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error GET: %v\n", err)
		return
	}
	defer resp.Body.Close()

	feed := Feed{}

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&feed)
	if err != nil {
		fmt.Printf("Error Decode: %v\n", err)
		return
	}

	for i, entry := range feed.Entry {
		fmt.Printf("%v. %v gave rating %v and wrote: %v on %v\n\n", i+1, entry.Author.Name, entry.Rating, entry.Title, entry.Updated)
	}
}