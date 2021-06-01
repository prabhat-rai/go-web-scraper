package handler

import (
	"context"
	"echoApp/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type (
	Handler struct {
		DB *mongo.Database
	}
)

func (h *Handler) authenticateUser(email string,password string) (user model.User, err error) {
	user, err = h.findUserByEmail(email)
	if err != nil {
		fmt.Printf("\n Find By Email : %v \n", err)
		return user,err
	}
	err = bcrypt.CompareHashAndPassword(user.Hashed_Password, []byte(password))
	return user,err
}

func (h *Handler) findUserByEmail(email string) (user model.User,err error){
	user = model.User{}
	// Find By Id
	userCollection := h.DB.Collection("users")
	dbContext := context.TODO()
	filter := bson.D{{"email", email}}
	err = userCollection.FindOne(dbContext, filter).Decode(&user)
	return user,err
}

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)
