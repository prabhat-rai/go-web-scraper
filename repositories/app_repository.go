package repositories

import (
	"context"
	"echoApp/model"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"go.mongodb.org/mongo-driver/bson"
)



type (

	AppRepository struct {
		DB *mongo.Database
	}
)




func (appsRepository *AppRepository) CreateApp(u *model.Apps) (err error){
	appCollection := appsRepository.DB.Collection("apps")
	dbContext := context.TODO()
	result, err := appCollection.InsertOne(dbContext, u)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Inserted Docs: ", result.InsertedID)
	return nil
}

func (appsRepository *AppRepository) UpdateAppStatus(u *model.Apps) (err error){
	filter := bson.D{{"_id", u.ID}}
	ctx := context.TODO()
	operation := "$set"
	appCollection := appsRepository.DB.Collection("apps")
	updateData := bson.M{operation: bson.M{"active": u.Active}}

	updateResult, err := appCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Updated Docs: ", updateResult)
	return nil
}