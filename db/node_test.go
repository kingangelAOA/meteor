package db

import (
	"fmt"
	"meteor/configs"
	"testing"
)

func Init() {
	m := configs.Mongo{
		User:     "kingangel",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     "27017",
		DB:       "meteor",
		Timeout:  5000,
	}
	InitMongo(m)
}
func TestNode(t *testing.T) {
	Init()
	ctx, cancel, c := getCollection("node")
	defer cancel()
	id, err := c.InsertOne(ctx, &Node{Name: "test", Type: "SCRIPT"})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(id)
	}
}
