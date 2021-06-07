package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Apps struct {
		ID        	primitive.ObjectID 	`json:"id" bson:"_id,omitempty"`
		Name      	string        		`json:"name" bson:"name"`
		GoogleAppId string        		`json:"google_app_id,omitempty" bson:"google_app_id"`
		IosAppId    string        		`json:"ios_app_id,omitempty" bson:"ios_app_id"`
		Active   	bool      			`json:"active,omitempty" bson:"active,omitempty"`
	}
)