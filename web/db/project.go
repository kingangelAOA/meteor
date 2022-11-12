package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Project struct
type Project struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Comments string             `json:"comments,omitempty" bson:"comments,omitempty"`
	Num      int                `json:"num" bson:"num,omitempty"`
}

type ProjectResponse struct {
	Projects []Project `json:"projects"`
	Total    int64     `json:"total"`
}

type NewProject struct {
	Name     string
	Comments string
}

type QueryProject struct {
	Page  int64  `form:"page"`
	Limit int64  `form:"limit"`
	Name  string `form:"name"`
}

func (qp *QueryProject) getSkip() int64 {
	return qp.Limit * (qp.Page - 1)
}

func (qp *QueryProject) GetOptions() *options.FindOptions {
	return options.Find().SetSkip(qp.getSkip()).SetLimit(qp.Limit)
}
