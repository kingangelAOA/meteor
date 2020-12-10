package core

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"testing"
)

type ExPaths struct {
	Paths openapi3.Paths
}

func (o *ExPaths)get() string {
	return "sssss"
}

func TestOpenApi(t *testing.T)  {
	s, _ := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("/Users/wenwen/work/meteor/core/test.yaml")
	j, _ := s.MarshalJSON()
	fmt.Println(string(j))
	//oldC := configs.Mongo{
	//	Host: "127.0.0.1",
	//	Port: "27017",
	//	Db:   "meteor",
	//	User: "admin",
	//	Password: "admin",
	//}
	//clientOld := db.MongoInit(oldC)
	//ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	//defer cancel()
	//if err := clientOld.Connect(ctx); err != nil {
	//	fmt.Println(err)
	//}
	//swaggers := clientOld.Database("meteor").Collection("swaggers")
	//cur, err := swaggers.Find(ctx, bson.D{})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//for cur.Next(ctx) {
	//	var swagger openapi3.Swagger
	//
	//	err := cur.Decode(&swagger)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(swagger.Paths)
	//}
}
