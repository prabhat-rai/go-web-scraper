package repositories

import (
	"context"
	"echoApp/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
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
		log.Fatal(err)
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