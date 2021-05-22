package handler

import (
	"context"
	"echoApp/model"
	"echoApp/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func (h *Handler) Home(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "index.tmpl", map[string]interface{}{
		"name": "World!",
	})
}

func (h *Handler) VerifyMongoDbQueries(c echo.Context) (err error) {
	userCollection := h.DB.Collection("users")
	dbContext := context.TODO()

	err = insertQueries(userCollection, dbContext)
	if err != nil {
		return services.HandleDbError(err)
	}

	err = fetchQueries(userCollection, dbContext)
	if err != nil {
		return services.HandleDbError(err)
	}

	err = updateQueries(userCollection, dbContext)
	if err != nil {
		return services.HandleDbError(err)
	}

	err = deleteQueries(userCollection, dbContext)
	if err != nil {
		return services.HandleDbError(err)
	}

	return c.JSON(http.StatusOK, "All Ok : Verified inserting, fetching, updating and deleting of records.")
}

func insertQueries(userCollection *mongo.Collection, dbContext context.Context) (err error){
	fmt.Println("Verifying Insert")

	// Insert 1 record
	objectId, _ := primitive.ObjectIDFromHex("60a550a5351458b8b460d762")
	record := &model.User{ID: objectId, Name: "Testing Four", Email: "test4@gmail.com", Phone: "+91444444444", Gender: "F"}
	result1, err := userCollection.InsertOne(dbContext, record)
	fmt.Println("Inserted Doc: ", result1.InsertedID)
	if err != nil {
		return err
	}

	// Insert Many
	insertRecords := []interface{}{
		&model.User{Name: "Testing One", Email: "test1@gmail.com", Phone: "+91111111111", Gender: "M"},
		&model.User{Name: "Testing Two", Email: "test2@gmail.com", Phone: "+91222222222", Gender: "F"},
		&model.User{Name: "Testing Three", Email: "test3@gmail.com", Phone: "+91333333333", Gender: "F"},
	}
	result, err := userCollection.InsertMany(dbContext, insertRecords)
	fmt.Println("Inserted Docs: ", result.InsertedIDs)
	if err != nil {
		return err
	}

	return nil
}

func fetchQueries(userCollection *mongo.Collection, dbContext context.Context) (err error) {
	fmt.Println("Verifying Fetch")
	var users []model.User
	user := model.User{}

	// Find By Id
	objectId, _ := primitive.ObjectIDFromHex("60a550a5351458b8b460d762")
	filter := bson.D{{"_id", objectId}}
	err = userCollection.FindOne(dbContext, filter).Decode(&user)
	if err != nil {
		return err
	}
	fmt.Printf("\n Find By Id : %v \n", user)


	// Find 1 record by email and hide _id from result
	filter = bson.D{{"email", "test1@gmail.com"}}
	opts := options.FindOne().SetProjection(bson.M{"ID": 0})
	err = userCollection.FindOne(dbContext, filter, opts).Decode(&user)
	if err != nil {
		return err
	}
	fmt.Printf("\n 1 Record By Email and Hide _id From Result : %v \n", user)


	// Get All matching Records
	filter = bson.D{{"gender", "M"}}
	findOpts := options.Find().SetProjection(bson.M{"ID": 0}).SetLimit(2)
	cursor, err := userCollection.Find(dbContext, filter, findOpts)
	if err != nil {
		return err
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&user)

		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}
	fmt.Printf("\n Get All matching Records : %v \n", users)

	return nil
}

func updateQueries(userCollection *mongo.Collection, dbContext context.Context) (err error){
	fmt.Println("Verifying Update")

	// Update By Id
	objectId, _ := primitive.ObjectIDFromHex("60a550a5351458b8b460d762")
	filter := bson.D{{"_id", objectId}}
	updateData := bson.D{
		{"$set",
			bson.D{
				{"phone", "+917777777777"},
				{"gender", "??"},
			},
		},
	}
	updateResult, err := userCollection.UpdateOne(dbContext, filter, updateData)
	if err != nil {
		return err
	}
	fmt.Printf("Updated documents: %+v\n", updateResult)

	// Update 1 Record
	condition := bson.D{{"email", "test1@gmail.com"}}
	updateData = bson.D{
		{"$set",
			bson.D{
				{"phone", "+9188887777"},
				{"gender", "F"},
			},
		},
	}

	updateResult, err = userCollection.UpdateOne(dbContext, condition, updateData)
	if err != nil {
		return err
	}
	fmt.Printf("Updated documents: %+v\n", updateResult)

	// Update Many Records
	pipeline := bson.M{
		"$or": []interface{}{
			bson.D{{"email", "test1@gmail.com"}},
			bson.D{{"email", "test2@gmail.com"}},
			bson.D{{"email", "test3@gmail.com"}},
			bson.D{{"phone", "+917777777777"}},
		},
	}

	updateData = bson.D{
		{"$set",
			bson.D{
				{"gender", ""},
			},
		},
	}

	updateResult, err = userCollection.UpdateMany(dbContext, pipeline, updateData)
	if err != nil {
		return err
	}
	fmt.Printf("Updated documents: %+v\n", updateResult)

	return nil
}

func deleteQueries(userCollection *mongo.Collection, dbContext context.Context) (err error){
	fmt.Println("Verifying Delete")

	// Delete Record
	filter := bson.D{{"email", "test3@gmail.com"}}
	deleteResult, err := userCollection.DeleteOne(dbContext, filter)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted %v documents \n", deleteResult.DeletedCount)


	// Delete all Matching records
	pipeline := bson.M{
		"$or": []interface{}{
			bson.D{{"email", "test1@gmail.com"}},
			bson.D{{"email", "test2@gmail.com"}},
			bson.D{{"email", "test3@gmail.com"}},
			bson.D{{"email", "test4@gmail.com"}},
		},
	}

	deleteResult, err = userCollection.DeleteMany(dbContext, pipeline)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted %v documents \n", deleteResult.DeletedCount)

	return nil
}