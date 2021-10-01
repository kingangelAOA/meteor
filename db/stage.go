package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Stage struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type string             `json:"type,omitempty" type:"name,omitempty"`
}
