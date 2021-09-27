package db

import (
	"fmt"
	"meteor/configs"
	"meteor/models"
	"testing"
)

func init() {
	m := configs.Mongo{
		User:     "kinganel",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     "27017",
		DB:       "meteor",
		Timeout:  5000,
	}
	InitMongo(m)
}

func TestNode(t *testing.T) {
	ctx, cancel, c := getCollection("node")
	defer cancel()
	id, err := c.InsertOne(ctx, &models.Node{Name: "test", Type: "SCRIPT"})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(id)
	}
}
