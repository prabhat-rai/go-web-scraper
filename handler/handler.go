package handler

import (
	"echoApp/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Handler struct {
		DB *mongo.Database
		UserRepository *repositories.UserRepository
	}
)
const (
	// Key (Should come from somewhere else).
	Key = "secret"
)
