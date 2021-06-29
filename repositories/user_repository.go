package repositories

import (
	"context"
	"echoApp/model"
	"echoApp/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math"
)

type (
	UserRepository struct {
		DB *mongo.Database
	}
)

func (userRepo *UserRepository) CreateUser(u *model.User) (err error) {
	// @TODO : VALIDATE REQUEST
	u.Hashed_Password, _ = bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	u.Password = ""
	userCollection := userRepo.DB.Collection("users")
	dbContext := context.TODO()
	result, err := userCollection.InsertOne(dbContext, u)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Inserted Docs: ", result.InsertedID)
	return nil
}

func (userRepo *UserRepository) AuthenticateUser(email string,password string) (user model.User, err error) {
	user, err = userRepo.FindUserByEmail(email)

	if err != nil {
		return user,err
	}

	err = bcrypt.CompareHashAndPassword(user.Hashed_Password, []byte(password))
	return user,err
}

func (userRepo *UserRepository) FindUserByEmail(email string) (user model.User,err error){
	user = model.User{}

	// Find By Id
	userCollection := userRepo.DB.Collection("users")
	dbContext := context.TODO()
	filter := bson.D{{"email", email}}
	err = userCollection.FindOne(dbContext, filter).Decode(&user)

	return user, err
}

func (userRepo *UserRepository) ListUsers(dataTableFilters *services.DataTableFilters,searchWord string) (userData model.UsersData,err error){
	userCollection := userRepo.DB.Collection("users")
	dbContext := context.TODO()

	findOptions := options.Find().SetLimit(dataTableFilters.Limit)

	findOptions.SetSort(bson.D{{dataTableFilters.SortColumnName, dataTableFilters.SortOrder}})
	findOptions.SetSkip(dataTableFilters.Offset)
	findOptions.SetLimit(dataTableFilters.Limit)

	count, err := userCollection.CountDocuments(dbContext, bson.D{})
//search word in name
	var searchFiltersinName = bson.D{}		
		searchFiltersinName = append(searchFiltersinName, bson.E{"name", primitive.Regex{Pattern: dataTableFilters.Search, Options: "i"}})

	cursor, err := userCollection.Find(dbContext, searchFiltersinName,findOptions)
	defer cursor.Close(dbContext)

	if err != nil {
		log.Println(err)
		return
	}
	for cursor.Next(dbContext) {
		user := model.User{}
		err := cursor.Decode(&user)

		if err != nil {
			log.Panic(err)
		}

		userData.Data = append(userData.Data, user)
	}

	userData.Total = count
	userData.LastPage = math.Ceil(float64(count / dataTableFilters.Limit))
	if userData.Data == nil {
		userData.Data = make([]model.User, 0)
	}
	return userData, err
}

func (userRepo *UserRepository) ChangePassword(email string,password string) (err error){
	// Find By Email
	userCollection := userRepo.DB.Collection("users")
	dbContext := context.TODO()
	hashedPassword,_ := bcrypt.GenerateFromPassword([]byte(password), 12)

	filter := bson.D{{"email", email}}
	update := bson.D{{"$set", bson.D{{"hpassword", hashedPassword}}}}

	result, err := userCollection.UpdateOne(dbContext, filter,update)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Updated records count",result.ModifiedCount)
	return
}

func (userRepo *UserRepository) UpdateUser(u *model.User) (err error){
	filter := bson.D{{"_id", u.ID}}
	ctx := context.TODO()
	operation := "$set"
	userCollection := userRepo.DB.Collection("users")
	updateData := bson.M{operation: bson.M{"name": u.Name, "email": u.Email, "role": u.Role, "phone": u.Phone}}
	updateResult, err := userCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Updated Docs: ", updateResult)
	return nil
}