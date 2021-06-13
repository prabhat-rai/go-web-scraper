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
	KeywordGroupRepository struct {
		DB *mongo.Database
	}

	AllKeywordGroups struct {
		Total    int64    				`json:"recordsTotal"`
		Filtered int64    				`json:"recordsFiltered"`
		Data     []model.KeywordGroup 	`json:"data"`
	}
)

func (keywordGroupRepo *KeywordGroupRepository) RetrieveKeywordGroups(dataTableFilters *services.DataTableFilters) (allKeywordGroups AllKeywordGroups) {
	finalSearchCondition := bson.D{}
	ctx := context.TODO()
	keywordCollection := keywordGroupRepo.DB.Collection("keyword_groups")

	var searchFilters bson.D

	if dataTableFilters.Search != "" {
		searchFilters = append(searchFilters, bson.E{"group_name", primitive.Regex{Pattern: dataTableFilters.Search, Options: ""}})
		//searchFilters = append(searchFilters, bson.E{"keywords", bson.M{"$in": dataTableFilters.Search}})
	}

	if len(searchFilters) > 0 {
		finalSearchCondition = bson.D{{"$or", []interface{}{
			searchFilters,
		}}}
	}

	// Set Find Options
	findOptions := options.Find().SetLimit(dataTableFilters.Limit)
	findOptions.SetSort(bson.D{{dataTableFilters.SortColumnName, dataTableFilters.SortOrder}})
	findOptions.SetSkip(dataTableFilters.Offset)

	count, err := keywordCollection.CountDocuments(ctx, bson.D{})
	cursor, err := keywordCollection.Find(ctx, finalSearchCondition, findOptions)
	filteredCount, err := keywordCollection.CountDocuments(ctx, finalSearchCondition)

	if err != nil {
		return allKeywordGroups
	}

	keywordGroup := model.KeywordGroup{}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&keywordGroup)

		if err != nil {
			log.Fatal(err)
		}

		allKeywordGroups.Data = append(allKeywordGroups.Data, keywordGroup)
	}

	allKeywordGroups.Total = count
	allKeywordGroups.Filtered = filteredCount

	if allKeywordGroups.Data == nil {
		allKeywordGroups.Data = make([]model.KeywordGroup, 0)
	}

	return allKeywordGroups
}