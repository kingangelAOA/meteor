package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"meteor/configs"
	"time"
)

const DataBase = "meteor"

var Mongo *MongoDB

type MongoDB struct {
	Client  *mongo.Client
	Timeout int64
}

//InitMongo *mongo.Client
func InitMongo(m configs.Mongo) {
	fmt.Println(m.GetMongoUrl())
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.GetMongoUrl()))
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
	p := Mongo.Client.Database(DataBase).Collection(c)
	return ctx, cancel, p
}
