package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	AppReview struct {
		ID        		primitive.ObjectID 	`json:"id" bson:"_id,omitempty"`
		ReviewId      	string        		`json:"review_id" bson:"review_id"`
		ReviewDate  	primitive.Timestamp	`json:"review_date" bson:"review_date"`
		UserName 		string		  		`json:"user_name" bson:"user_name"`
		Title			string				`json:"review_title" bson:"review_title"`
		Description  	string        		`json:"review_description" bson:"review_description"`
		Rating    		int        			`json:"rating" bson:"rating"`
		CreatedAt 		primitive.Timestamp	`json:"created_at" bson:"created_at"`
		UpdatedAt 		primitive.Timestamp	`json:"updated_at" bson:"updated_at"`
		Platform 		string				`json:"platform" bson:"platform"`
		Version 		string				`json:"version" bson:"version"`
		Concept 		string				`json:"concept" bson:"concept"`
		Keywords   		[]string      		`json:"keywords" bson:"keywords"`
	}
)
