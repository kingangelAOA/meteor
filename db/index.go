package db

import (
	"context"
	"fmt"
	"meteor/configs"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DataBase = "meteor"

var Mongo *MongoDB

type MongoDB struct {
	Client  *mongo.Client
	Timeout int64
}

//InitMongo *mongo.Client
func InitMongo(m configs.Mongo) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.GetMongoUrl()))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	Mongo = &MongoDB{
		Client:  client,
		Timeout: m.Timeout,
	}
}

func CreateMongo(m configs.Mongo) *MongoDB {
	fmt.Println(m.GetMongoUrl())
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.GetMongoUrl()))
	if err != nil {
		panic(err)
	}
	return &MongoDB{
		Client:  client,
		Timeout: m.Timeout,
	}
}

func (m *MongoDB) getConnectCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), time.Duration(m.Timeout)*time.Millisecond)
}

func getCollection(c string) (context.Context, context.CancelFunc, *mongo.Collection) {
	ctx, cancel := Mongo.getConnectCtx()
	col := Mongo.Client.Database(DataBase).Collection(c)
	return ctx, cancel, col
}

type DBRef struct {
	Ref string             `bson:"$ref"`
	ID  primitive.ObjectID `bson:"$id"`
	DB  string             `bson:"$db"`
}

func NewDBRef(ref string, id primitive.ObjectID) DBRef {
	return DBRef{
		Ref: ref,
		ID:  id,
		DB:  DataBase,
	}
}
