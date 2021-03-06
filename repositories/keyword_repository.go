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
	KeywordRepository struct {
		DB *mongo.Database
	}

	AllKeywords struct {
		Total    int64    			`json:"recordsTotal"`
		Filtered int64    			`json:"recordsFiltered"`
		Data     []model.Keyword 	`json:"data"`
	}
)


func (keywordRepo *KeywordRepository) RetrieveKeywords(dataTableFilters *services.DataTableFilters) (allKeywords AllKeywords) {
	var keyword model.Keyword
	var searchFilters bson.D

	finalSearchCondition := bson.D{}
	ctx := context.TODO()
	keywordCollection := keywordRepo.DB.Collection("keywords")

	if dataTableFilters.Search != "" {
		searchFilters = append(searchFilters, bson.E{"name", primitive.Regex{Pattern: dataTableFilters.Search, Options: ""}})
	}

	if len(searchFilters) > 0 {
		finalSearchCondition = bson.D{{"$or", []interface{}{
			searchFilters,
		}}}
	}

	// Set Find Options
	findOptions := options.Find()

	if dataTableFilters.Limit != 0 {
		findOptions.SetLimit(dataTableFilters.Limit)
	}

	if dataTableFilters.SortColumnName != "" {
		findOptions.SetSort(bson.D{{dataTableFilters.SortColumnName, dataTableFilters.SortOrder}})
	}

	findOptions.SetSkip(dataTableFilters.Offset)

	count, err := keywordCollection.CountDocuments(ctx, bson.D{})
	cursor, err := keywordCollection.Find(ctx, finalSearchCondition, findOptions)
	filteredCount, err := keywordCollection.CountDocuments(ctx, finalSearchCondition)

	if err != nil {
		return allKeywords
	}

	for cursor.Next(context.TODO()) {
		keyword = model.Keyword{}
		err := cursor.Decode(&keyword)

		if err != nil {
			log.Panic(err)
		}

		allKeywords.Data = append(allKeywords.Data, keyword)
	}

	allKeywords.Total = count
	allKeywords.Filtered = filteredCount

	if allKeywords.Data == nil {
		allKeywords.Data = make([]model.Keyword, 0)
	}
	return allKeywords
}
func (keywordRepo *KeywordRepository) CreateKeyword(u *model.Keyword) (err error) {
	keywordCollection := keywordRepo.DB.Collection("keywords")
	dbContext := context.TODO()
	result, err := keywordCollection.InsertOne(dbContext, u)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Inserted Docs: ", result.InsertedID)
	return nil
}

func (keywordRepo *KeywordRepository) UpdateActiveStatus(u *model.Keyword) (err error){
	filter := bson.D{{"_id", u.ID}}
	ctx := context.TODO()
	operation := "$set"
	keywordCollection := keywordRepo.DB.Collection("keywords")
	updateData := bson.M{operation: bson.M{"active": u.Active}}
	updateResult, err := keywordCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Updated Docs: ", updateResult)
	return nil
}



