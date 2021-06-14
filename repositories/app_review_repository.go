package repositories

import (
	"context"
	"echoApp/model"
	"echoApp/services"
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
	result, err := appReviewCollection.InsertMany(dbContext, insertRecords)

	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Inserted Docs: ", result.InsertedIDs)
	return nil
}

func (appReviewRepo *AppReviewRepository) RetrieveBulkReviews(dataTableFilters *services.DataTableFilters, filters map[string] string, keywords []string) (allReviews AllReviews) {
	finalSearchCondition := bson.D{}
	ctx := context.TODO()
	appReviewCollection := appReviewRepo.DB.Collection("app_reviews")

	var andFilters bson.D
	var searchFilters bson.D
	var keywordFilters bson.M

	for key, value := range filters {
		andFilters = append(andFilters, bson.E{key, value})
	}

	if len(keywords) > 0 {
		keywordFilters = bson.M{"keywords": bson.M{"$in": keywords}}
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

	if len(keywordFilters) == 0 {
		keywordFilters = bson.M{}
	}

	if len(andFilters) > 0 || len(searchFilters) > 0 {
		finalSearchCondition = bson.D{
			{ "$and", []interface{}{
				keywordFilters,
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

func (appReviewRepo *AppReviewRepository) GetLatestReviewId(platform string, concept string) string {
	ctx := context.TODO()
	appReviewCollection := appReviewRepo.DB.Collection("app_reviews")
	review := model.AppReview{}

	filter := bson.D{{"platform", platform}, {"concept", concept}}

	opts := options.FindOne().SetSort(bson.D{{"review_date", -1}})
	err := appReviewCollection.FindOne(ctx, filter, opts).Decode(&review)
	if err != nil {
		log.Println(err)
		return ""
	}

	return review.ReviewId
}