package handler

import (
	"echoApp/model"
	"echoApp/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func (h *Handler) Home(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "index.tmpl", map[string]interface{}{
		"name": "World!",
	})
}

func (h *Handler) VerifyMongoDbQueries(c echo.Context) (err error) {
	// Connect to DB
	db := h.DB.Clone()
	defer db.Close()
	userCollection := db.DB(h.DbName).C("users")

	err = insertQueries(userCollection)
	if err != nil {
		return services.HandleDbError(err)
	}

	err = fetchQueries(userCollection)
	if err != nil {
		return services.HandleDbError(err)
	}

	err = updateQueries(userCollection)
	if err != nil {
		return services.HandleDbError(err)
	}

	err = deleteQueries(userCollection)
	if err != nil {
		return services.HandleDbError(err)
	}

	return c.JSON(http.StatusOK, "All Ok : Verified inserting, fetching, updating and deleting of records.")
}

func insertQueries(userCollection *mgo.Collection) (err error){
	fmt.Println("Verifying Insert")
	// Insert Record
	err = userCollection.Insert(
		&model.User{Name: "Testing One", Email: "test1@gmail.com", Phone: "+91111111111", Gender: "M"},
		&model.User{Name: "Testing Two", Email: "test2@gmail.com", Phone: "+91222222222", Gender: "F"},
		&model.User{Name: "Testing Three", Email: "test3@gmail.com", Phone: "+91333333333", Gender: "F"},
		&model.User{ID: bson.ObjectIdHex("60a550a5351458b8b460d762"), Name: "Testing Four", Email: "test4@gmail.com", Phone: "+91444444444", Gender: "F"},
	)

	if err != nil {
		return err
	}

	return nil
}

func fetchQueries(userCollection *mgo.Collection) (err error) {
	fmt.Println("Verifying Fetch")
	var users []model.User
	user := model.User{}

	// Find By Id
	err = userCollection.FindId(bson.ObjectIdHex("60a550a5351458b8b460d762")).One(&user)
	if err != nil {
		return err
	}
	fmt.Printf("\n Find By Id : %v \n", user)


	// Find 1 record by email and hide _id from result
	err = userCollection.Find(bson.M{"email" : "test1@gmail.com"}).Select(bson.M{"_id" : 0}).One(&user)
	if err != nil {
		return err
	}
	fmt.Printf("\n 1 Record By Email and Hide _id From Result : %v \n", user)


	// Get All matching Records
	err = userCollection.Find(bson.M{"gender" : "M"}).Select(bson.M{"password" : 0}).All(&users)
	if err != nil {
		return err
	}
	fmt.Printf("\n Get All matching Records : %v \n", users)

	return nil
}

func updateQueries(userCollection *mgo.Collection) (err error){
	fmt.Println("Verifying Update")

	// Update Record
	condition := bson.M{"email": "test1@gmail.com"}
	updateData := bson.M{"$set": bson.M{"phone": "+9188887777", "gender" : "F"}}
	err = userCollection.Update(condition, updateData)
	if err != nil {
		return err
	}

	return nil
}

func deleteQueries(userCollection *mgo.Collection) (err error){
	fmt.Println("Verifying Delete")

	// Delete Record
	err = userCollection.Remove(bson.M{"email" : "test3@gmail.com"})
	if err != nil {
		return err
	}

	// Delete all Matching records

	// @note : We can either use bson.D or bson.M to apply logical operators
	//pipeline := bson.D{
	//	{"$or", []interface{}{
	//		bson.D{{"email", "test1@gmail.com"}},
	//		bson.D{{"email", "test2@gmail.com"}},
	//	}},
	//}

	pipeline := bson.M{
		"$or": []interface{}{
			bson.M{"email": "test1@gmail.com"},
			bson.M{"email": "test2@gmail.com"},
			bson.M{"email": "test3@gmail.com"},
			bson.M{"phone": "+91444444444"},
		},
	}

	_, err = userCollection.RemoveAll(pipeline)
	if err != nil {
		return err
	}

	return nil
}