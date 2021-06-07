package repositories

import (
	"context"
	"echoApp/model"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AppReviewRepository struct {
		DB *mongo.Database
	}
)

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