package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Operation struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	StagId       primitive.ObjectID   `json:"stagId,omitempty" bson:"stagId,omitempty"`
	Performances []primitive.ObjectID `json:"performances,omitempty" bson:"performances,omitempty"`
	Assert       Assert
}

type Performance struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RT99      int                `json:"RT99,omitempty" bson:"RT99,omitempty"`
	RT95      int                `json:"RT95,omitempty" bson:"RT95,omitempty"`
	RT90      int                `json:"RT90,omitempty" bson:"RT90,omitempty"`
	AverageRT int                `json:"averagRT,omitempty" bson:"averagRT,omitempty"`
	ErrRate   float64            `json:"errRate,omitempty" bson:"errRate,omitempty"`
	QPS       int                `json:"QPS,omitempty" bson:"QPS,omitempty"`
}

type Assert struct {
}
