package services

import (
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

func FetchReview(platform string) {
	if platform == "ios" {
		fetchAppStoreReviews()
	} else {
		fetchPlayStoreReviews()
	}

}

func fetchAppStoreReviews () {
	url := "https://itunes.apple.com/ae/rss/customerreviews/id=1161310479/page=1/sortBy=mostRecent/xml";
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

func fetchPlayStoreReviews() {
	r := reviews.New("com.landmarkgroup.maxstores", reviews.Options{
		Number: 100,
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