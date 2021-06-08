package repositories

import (
	"context"
	"echoApp/model"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	AppReviewRepository struct {
		DB *mongo.Database
	}
)


type AppReviewData struct {
	
	ReviewId      	string        		`json:"review_id" bson:"review_id"`
	ReviewDate  	string        		`json:"review_date" bson:"review_date"`
	UserName 		string		  		`json:"user_name" bson:"user_name"`
	Title			string				`json:"review_title" bson:"review_title"`
	Description  	string        		`json:"review_description" bson:"review_description"`
	Rating    		string        		`json:"rating" bson:"rating"`
	CreatedAt 		string				`json:"created_at" bson:"created_at"`
	UpdatedAt 		string				`json:"updated_at" bson:"updated_at"`
	Platform 		string				`json:"platform" bson:"platform"`
	Version 		string				`json:"version" bson:"version"`
	Concept 		string				`json:"concept" bson:"concept"`
	Keywords   		[]string      		`json:"keywords" bson:"keywords,omitempty"`
}

type AllReviews struct {
	AppReview []model.AppReview `json:"app_reviews"`
}
	 

func (appReviewRepo *AppReviewRepository) AddBulkReviews(appReviews []*model.AppReview) (err error) {
	var insertRecords []interface{}
	for _, elem := range appReviews {
		insertRecords = append(insertRecords, elem)
	}

	appReviewCollection := appReviewRepo.DB.Collection("app_reviews")
	dbContext := context.TODO()
	//insertRecords := []interface{}{appReviews}
	result, err := appReviewCollection.InsertMany(dbContext, insertRecords)

	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	fmt.Println("Inserted Docs: ", result.InsertedIDs)
	return nil
}

func (appReviewRepo *AppReviewRepository) RetrieveBulkReviews() (allReviews AllReviews) {
	findOptions := options.Find()
	findOptions.SetLimit(2)
	// var results []*Review
	ctx := context.TODO()
	review := model.AppReview{}
	appReviewCollection := appReviewRepo.DB.Collection("app_reviews")
	findOpts := options.Find().SetProjection(bson.M{"ID": 0}).SetLimit(10)
	cursor, err := appReviewCollection.Find(ctx, bson.D{}, findOpts)
	if err != nil {
		return allReviews
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&review)

		if err != nil {
			log.Fatal(err) 
		}

		allReviews.AppReview = append(allReviews.AppReview, review)
	}
	return allReviews
	
}