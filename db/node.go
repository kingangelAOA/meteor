package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Node struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Type   string             `json:"type,omitempty" bson:"type,omitempty"`
	Config map[string]string  `json:"config,omitempty" bson:"config,omitempty"`
}

type Http struct {
	Name   string `json:"name,omitempty" bson:"name,omitempty"`
	Method string `json:"method,omitempty" bson:"method,omitempty"`
	Host   string `json:"host,omitempty" bson:"host,omitempty"`
	Port   string `json:"port,omitempty" bson:"port,omitempty"`
	Path   string `json:"path,omitempty" bson:"path,omitempty"`
	Body   string `json:"body,omitempty" bson:"body,omitempty"`
	Header []KV   `json:"header,omitempty" bson:"header,omitempty"`
	Query  []KV   `json:"query,omitempty" bson:"query,omitempty"`
}

type Script struct {
	Key  primitive.ObjectID `json:"key,omitempty" bson:"key,omitempty"`
	Type string             `json:"type,omitempty" bson:"type,omitempty"`
	Code string             `json:"code,omitempty" bson:"code,omitempty"`
}

type KV struct {
	Key   string `json:"key,omitempty" bson:"key,omitempty"`
	Value string `json:"value,omitempty" bson:"value,omitempty"`
}
