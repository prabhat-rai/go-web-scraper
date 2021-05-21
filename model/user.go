package model

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	User struct {
		ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Name      string        `json:"name" bson:"name"`
		Email     string        `json:"email" bson:"email"`
		Phone  	  string        `json:"phone,omitempty" bson:"phone"`
		Gender    string        `json:"gender,omitempty" bson:"gender"`
		Hobbies   []string      `json:"hobbies,omitempty" bson:"hobbies,omitempty"`
	}
)
