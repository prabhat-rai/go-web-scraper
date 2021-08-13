package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	KeywordGroup struct {
		ID        		primitive.ObjectID 	`json:"id" bson:"_id,omitempty"`
		GroupName 		string		  		`json:"group_name" bson:"group_name"`
		Keywords   		[]string      		`json:"keywords" bson:"keywords"`
		Active   		bool      			`json:"active,omitempty" bson:"active,omitempty"`
		Subscribers 	[]string      		`json:"subscribers,omitempty" bson:"subscribers,omitempty"`
		SubscribeAction string      		`json:"subscribe_action"`
	}
)
