package handler

import (
	"context"
	"echoApp/model"
	"echoApp/services"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (h *Handler) Home(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "index.tmpl", map[string]interface{}{
		"name": "Admin",
	})
}

func (h *Handler) Login(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "login.tmpl", map[string]interface{}{
		"name": "Admin",
	})
}

func (h *Handler) LoginForm(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "login.page.tmpl", map[string]interface{}{
	})
}

func (h *Handler) Login(c echo.Context) error {
	u := new(model.User)
	err := c.Bind(u);
	if err != nil {
		//place holder to render login with error message
		return c.Render(http.StatusOK, "login.page.tmpl", map[string]interface{}{
		})
	}
	user,err :=h.authenticateUser(u.Email,u.Password)
	fmt.Printf("User authenticated",user.Email)
	if err != nil {
		updateSession(c)
	}
	//placeholder to redirect to dashboard
	http.Redirect(c.Response(), c.Request(), "/dashboard", http.StatusSeeOther)
	return nil
}

func updateSession(c echo.Context) {
	sess, _ := session.Get("session", c)
	sess.Values["authenticated"] = true
	sess.Save(c.Request(),c.Response())
}

func (h *Handler) RegisterForm(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "register.page.tmpl", map[string]interface{}{
	})
}

func (h *Handler) Register(c echo.Context) (err error) {
	u := new(model.User)
	err = c.Bind(u);
	if err != nil {
		//place holder to render register page with error message
		return c.Render(http.StatusOK, "register.page.tmpl", map[string]interface{}{
		})
	}
	err = h.ceateUser(u)
	if err != nil {
		updateSession(c)
	}
	return err
}

func (h *Handler) ceateUser(u *model.User) (err error) {
	u.Hashed_Password, _ = bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	u.Password = ""
	userCollection := h.DB.Collection("users")
	dbContext := context.TODO()
	result, err := userCollection.InsertOne(dbContext, u)
	fmt.Println("Inserted Docs: ", result.InsertedID)
	return err
}

func (h *Handler) Logout(c echo.Context) (err error) {
	sess, _ := session.Get("session", c)
	sess.Values["authenticated"] = false
	sess.Save(c.Request(), c.Response())

	return err
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