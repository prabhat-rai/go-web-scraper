package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Keyword struct {
		ID        	primitive.ObjectID 	`json:"id" bson:"_id,omitempty"`
		Name 		string		  		`json:"name" bson:"name"`
		Active   	bool      			`json:"active,omitempty" bson:"active,omitempty"`
	}
)
