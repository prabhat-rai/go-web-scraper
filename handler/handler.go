package handler

import (
	"gopkg.in/mgo.v2"
)

type (
	Handler struct {
		DB *mgo.Session
		DbName string
	}
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)
