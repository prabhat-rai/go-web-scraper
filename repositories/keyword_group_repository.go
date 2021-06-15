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

func (keywordGroupRepo *KeywordGroupRepository) RetrieveKeywordGroups(dataTableFilters *services.DataTableFilters, userData *model.User, activeRecords string) (allKeywordGroups AllKeywordGroups) {
	finalSearchCondition := bson.D{}
	ctx := context.TODO()
	keywordCollection := keywordGroupRepo.DB.Collection("keyword_groups")

	var searchFilters bson.D
	var andFilters bson.D

	if dataTableFilters.Search != "" {
		searchFilters = append(searchFilters, bson.E{"group_name", primitive.Regex{Pattern: dataTableFilters.Search, Options: ""}})
		//searchFilters = append(searchFilters, bson.E{"keywords", bson.M{"$in": dataTableFilters.Search}})
	}

	if activeRecords == "1" {
		andFilters = append(andFilters, bson.E{"active", true})
	}

	if len(andFilters) == 0 {
		andFilters = bson.D{}
	}

	if len(searchFilters) == 0 {
		searchFilters = bson.D{}
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
		return allKeywordGroups
	}

	keywordGroup := model.KeywordGroup{}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&keywordGroup)

		if err != nil {
			log.Fatal(err)
		}

		if services.InArray(userData.Email, keywordGroup.Subscribers) == true {
			keywordGroup.SubscribeAction = "<a href='javascript:void(0)' " +
				"class='btn btn-danger btn-icon-split btn-sm' " +
				"onclick='webScrapperApp.changeKeyGroupSubscription(\""+ keywordGroup.ID.Hex() +"\", 0)'>" +
				"<span class='text'>Unsubscribe</span> " +
				"<span class='icon text-white-50'><i class='fas fa-times'></i> </span> </a>"
		} else {
			keywordGroup.SubscribeAction = "<a href='javascript:void(0)' class='btn btn-primary btn-icon-split btn-sm'" +
				"onclick='webScrapperApp.changeKeyGroupSubscription(\""+ keywordGroup.ID.Hex() +"\", 1)'>" +
				"<span class='text'>Subscribe</span> " +
				"<span class='icon text-white-50'><i class='fas fa-rss'></i> </span> </a>"
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

func (keywordGroupRepo *KeywordGroupRepository) GetKeywordsForGroup(groupId string) []string{
	keywordGroup := model.KeywordGroup{}
	objectId, _ := primitive.ObjectIDFromHex(groupId)
	filter := bson.D{{"_id", objectId}}
	ctx := context.TODO()
	keywordCollection := keywordGroupRepo.DB.Collection("keyword_groups")
	err := keywordCollection.FindOne(ctx, filter).Decode(&keywordGroup)

	if err != nil {
		log.Fatal(err)
		return make([]string, 0)
	}

	return keywordGroup.Keywords
}

func (keywordGroupRepo *KeywordGroupRepository) UpdateSubscriptionForUser(keyGroupId string, subscriptionStatus string, userEmail string) int64{
	objectId, _ := primitive.ObjectIDFromHex(keyGroupId)
	filter := bson.D{{"_id", objectId}}
	ctx := context.TODO()
	keywordCollection := keywordGroupRepo.DB.Collection("keyword_groups")

	operation := "$push"

	if subscriptionStatus == "0" {
		operation = "$pull"
	}

	updateData := bson.M{operation: bson.M{"subscribers": userEmail}}

	_, err := keywordCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	return 1
}


func (keywordGroupRepo *KeywordGroupRepository) GetGroupsWithActiveSubscribers() (allKeywordGroups []model.KeywordGroup) {
	keywordCollection := keywordGroupRepo.DB.Collection("keyword_groups")
	ctx := context.TODO()

	filter := bson.D{{"$or", []interface{}{
		bson.D{{"subscribers.0", bson.M{"$exists": true}}},
	}}}

	findOptions := options.Find()
	cursor, err := keywordCollection.Find(ctx, filter, findOptions)

	keywordGroup := model.KeywordGroup{}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&keywordGroup)

		if err != nil {
			log.Fatal(err)
		}

		allKeywordGroups = append(allKeywordGroups, keywordGroup)
	}


	if err != nil {
		log.Println(err)
		return
	}

	return allKeywordGroups
}
