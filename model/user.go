package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	User struct {
		ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name      string        `json:"name" bson:"name"`
		Email     string        `json:"email" bson:"email"`
		Password string			`json:"password,omitempty" bson:"password,omitempty"`
		Hashed_Password []byte	`bson:"hpassword"`
		Phone  	  string        `json:"phone,omitempty" bson:"phone"`
		Gender    string        `json:"gender,omitempty" bson:"gender"`
		Hobbies   []string      `json:"hobbies,omitempty" bson:"hobbies,omitempty"`
	}
)
