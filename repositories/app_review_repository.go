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
	"time"
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

func (appReviewRepo *AppReviewRepository) AddBulkReviews(appReviews []*model.AppReview) (insertedIds interface {}, err error) {
	var insertRecords []interface{}
	for _, elem := range appReviews {
		insertRecords = append(insertRecords, elem)
	}

	appReviewCollection := appReviewRepo.DB.Collection("app_reviews")
	dbContext := context.TODO()
	result, err := appReviewCollection.InsertMany(dbContext, insertRecords)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result.InsertedIDs, nil
}

func (appReviewRepo *AppReviewRepository) RetrieveBulkReviews(dataTableFilters *services.DataTableFilters, filters map[string] string, keywords []string) (allReviews AllReviews) {
	var review model.AppReview
	var andFilters bson.D
	var searchFilters bson.D
	var keywordFilters bson.M

	finalSearchCondition := bson.D{}
	ctx := context.TODO()
	appReviewCollection := appReviewRepo.DB.Collection("app_reviews")

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

	for cursor.Next(context.TODO()) {
		review = model.AppReview{}
		err := cursor.Decode(&review)

		if err != nil {
			log.Panic(err)
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

func (appReviewRepo *AppReviewRepository) GetReviewsWithMatchingKeywords(keywords []string, objectIds interface{}) (allReviews []model.AppReview) {
	var review model.AppReview
	finalSearchCondition := bson.D{}
	var andFilters bson.M
	//var searchFilters bson.D
	var keywordFilters bson.M

	if len(keywords) > 0 {
		keywordFilters = bson.M{"keywords": bson.M{"$in": keywords}}
	}

	if objectIds != nil {
		andFilters = bson.M{"_id": bson.M{"$in": objectIds}}
	}

	ctx := context.TODO()
	appReviewCollection := appReviewRepo.DB.Collection("app_reviews")

	// Set Find Options
	findOptions := options.Find()

	finalSearchCondition = bson.D{
		{"$and", []interface{}{
			keywordFilters,
			andFilters,
			//bson.D{{"$or", []interface{}{
			//	searchFilters,
			//}}},
		}},
	}

	cursor, err := appReviewCollection.Find(ctx, finalSearchCondition, findOptions)

	if err != nil {
		return allReviews
	}

	for cursor.Next(context.TODO()) {
		review = model.AppReview{}
		err := cursor.Decode(&review)

		if err != nil {
			log.Panic(err)
		}

		allReviews = append(allReviews, review)
	}

	return allReviews
}

//group by count
func (appReviewRepo *AppReviewRepository) CountReviews(groupByAttribute string, difference int, differenceIn string) []bson.M {
	var workingTime primitive.Timestamp
	dbContext := context.TODO()
	reviewCollection := appReviewRepo.DB.Collection("app_reviews")

	switch differenceIn {
		case "months":
			workingTime = primitive.Timestamp{T:uint32(time.Now().AddDate(0,-difference,0).Unix())}
		case "years":
			workingTime = primitive.Timestamp{T:uint32(time.Now().AddDate(-difference,0,0).Unix())}
		case "days":
			workingTime = primitive.Timestamp{T:uint32(time.Now().AddDate(0,0,-difference).Unix())}
		default:
			fmt.Println("Unknown difference type encountered, switching to days")
			workingTime = primitive.Timestamp{T:uint32(time.Now().AddDate(0,0,-difference).Unix())}
	}
	matchStage := bson.D{
		{"$match", bson.D{
			{"review_date", bson.D{{"$gt",workingTime}}},
		}},
	}

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$" + groupByAttribute},
			{"count", bson.D{{"$sum", 1}}},
		}},
	}

	cur,err := reviewCollection.Aggregate(dbContext,mongo.Pipeline{matchStage,groupStage})
	var reviewAggregator []bson.M

	if cur == nil {
		fmt.Println("cursor nil")
		return reviewAggregator
	}

	if err = cur.All(dbContext, &reviewAggregator); err != nil {
		log.Panic(err)
	}
	return reviewAggregator
}

func (appReviewRepo *AppReviewRepository) DateWiseReviews(groupByAttribute string, differenceInDays int, datatype string) []bson.M {
	var workingTime primitive.Timestamp
	var reviewDateGroup bson.M
	dbContext := context.TODO()
	reviewCollection := appReviewRepo.DB.Collection("app_reviews")
	workingTime = primitive.Timestamp{T:uint32(time.Now().AddDate(0,0,-differenceInDays).Unix())}

	reviewDateGroup = bson.M{
		"$dateToString": bson.D{
			{"format", "%Y-%m-%d"},
			{"date", "$review_date"},
		},
	}

	matchStage := bson.D{
		{"$match", bson.D{
			{"review_date", bson.D{{"$gt",workingTime}}},
		}},
	}
	if datatype == "review"{
		groupStage := bson.D{{
			"$group", bson.D{
				{"_id", bson.D{
					{groupByAttribute, "$" + groupByAttribute},
					{"review_date", reviewDateGroup},
				},
				}, {"count", bson.D{
					{"$sum", 1},
				}},
			},
		}}
	
		cur,err := reviewCollection.Aggregate(dbContext,mongo.Pipeline{matchStage,groupStage})
		var reviewAggregator []bson.M
		if cur == nil {
			fmt.Println("cursor nil")
			return reviewAggregator
		}
	
		if err = cur.All(dbContext, &reviewAggregator); err != nil {
			log.Panic(err)
		}
		return reviewAggregator
	} else{
		groupStage := bson.D{{
			"$group", bson.D{
				{"_id", bson.D{
					{groupByAttribute, "$" + groupByAttribute},
					{"review_date", reviewDateGroup},
				},
				}, {"count", bson.D{
					{"$avg", "$rating"},
				}},
			},
		}}
	
		cur,err := reviewCollection.Aggregate(dbContext,mongo.Pipeline{matchStage,groupStage})
		var reviewAggregator []bson.M
		if cur == nil {
			fmt.Println("cursor nil")
			return reviewAggregator
		}
	
		if err = cur.All(dbContext, &reviewAggregator); err != nil {
			log.Panic(err)
		}
		return reviewAggregator
	}
	
}