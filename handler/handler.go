package handler

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Handler struct {
		DB *mongo.Database
	}
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)
