package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Node struct {
	ID     primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string                 `json:"name,omitempty" bson:"name,omitempty"`
	Type   string                 `json:"type,omitempty" bson:"type,omitempty"`
	Config map[string]interface{} `json:"config,omitempty" bson:"config,omitempty"`
}
