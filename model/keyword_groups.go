package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	KeywordGroup struct {
		ID        	primitive.ObjectID 	`json:"id" bson:"_id,omitempty"`
		GroupName 	string		  		`json:"group_name" bson:"group_name"`
		Keywords   	[]string      		`json:"keywords" bson:"keywords,omitempty"`
		Active   	bool      			`json:"active,omitempty" bson:"active,omitempty"`
	}
)
