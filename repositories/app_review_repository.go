package repositories

import (
	"context"
	"echoApp/model"
	"echoApp/services"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type (
	AppReviewRepository struct {
		DB *mongo.Database
	}

	AllReviews struct {
		Total    int64    `json:"recordsTotal"`
		Filtered int64    `json:"recordsFiltered"`
		Data     []model.AppReview `json:"data"`
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

func (appReviewRepo *AppReviewRepository) RetrieveBulkReviews(dataTableFilters *services.DataTableFilters, filters map[string] string) (allReviews AllReviews) {
	finalSearchCondition := bson.D{}
	ctx := context.TODO()
	appReviewCollection := appReviewRepo.DB.Collection("app_reviews")

	var andFilters bson.D
	var searchFilters bson.D

	for key, value := range filters {
		andFilters = append(andFilters, bson.E{key, value})
	}

	if dataTableFilters.Search != "" {
		searchFilters = append(searchFilters, bson.E{"review_title", primitive.Regex{Pattern: dataTableFilters.Search, Options: ""}})
		searchFilters = append(searchFilters, bson.E{"review_description", primitive.Regex{Pattern: dataTableFilters.Search, Options: ""}})
	}

	if len(searchFilters) == 0 {
		searchFilters = bson.D{}
	}

	if len(andFilters) == 0 {
		andFilters = bson.D{}
	}

	if len(andFilters) > 0 || len(searchFilters) > 0 {
		finalSearchCondition = bson.D{
			{ "$and", []interface{}{
				andFilters,
				bson.D{{"$or", []interface{}{
					searchFilters,
				}}},
			}},
		}
	}

	// Set Find Options
	findOptions := options.Find().SetLimit(dataTableFilters.Limit)
	findOptions.SetSort(bson.D{{dataTableFilters.SortColumnName, dataTableFilters.SortOrder}})
	findOptions.SetSkip(dataTableFilters.Offset)

	count, err := appReviewCollection.CountDocuments(ctx, bson.D{})
	cursor, err := appReviewCollection.Find(ctx, finalSearchCondition, findOptions)
	filteredCount, err := appReviewCollection.CountDocuments(ctx, finalSearchCondition)

	if err != nil {
		return allReviews
	}

	review := model.AppReview{}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&review)

		if err != nil {
			log.Fatal(err) 
		}

		allReviews.Data = append(allReviews.Data, review)
	}

	allReviews.Total = count
	allReviews.Filtered = filteredCount

	if allReviews.Data == nil {
		allReviews.Data = make([]model.AppReview, 0)
	}

	return allReviews
	
}

//group by count
func (appReviewRepo *AppReviewRepository) CountReviews(collection string, groupbyattr string) []bson.M {
	reviewCollection := appReviewRepo.DB.Collection(collection)
	dbContext := context.TODO()
	groupStage := bson.D{{"$group", bson.D{{"_id", groupbyattr}, {"count", bson.D{{"$sum", 1}}}}}}

	cur,err := reviewCollection.Aggregate(dbContext,mongo.Pipeline{groupStage})
	var showsWithInfo []bson.M
	if err = cur.All(dbContext, &showsWithInfo); err != nil {
		panic(err)
	}
	return showsWithInfo
}