package models

import (
	"github.com/getkin/kin-openapi/openapi3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreateUpdateSwaggerRequest update swagger request
type CreateUpdateSwaggerRequest struct {
	ID      primitive.ObjectID `bson:"_id"`
	Swagger string             `bson:"swagger"`
}

//GetSwagger get swagger object
func (cusr *CreateUpdateSwaggerRequest) GetSwagger() (*openapi3.Swagger, error) {
	return openapi3.NewSwaggerLoader().LoadSwaggerFromData([]byte(cusr.Swagger))
}

//QuerySwaggers query swagger
type QuerySwaggers struct {
	ProjectID string `form:"projectID,omitempty"`
	Version   string `form:"version,omitempty"`
}

type SwaggerFromDB struct {
	Version string `bson:"version"`
	Swagger string `bson:"swagger"`
}
