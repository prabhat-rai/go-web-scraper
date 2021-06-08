package handler

import (
	"echoApp/conf"
	"echoApp/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Handler struct {
		DB *mongo.Database
		UserRepository *repositories.UserRepository
		AppReviewRepository *repositories.AppReviewRepository
		Config *conf.Config
	}
)
const (
	// Key (Should come from somewhere else).
	Key = "secret"
)
