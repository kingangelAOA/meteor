package db

import (
	"fmt"
	"meteor/configs"
	"meteor/models"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Init() {
	m := configs.Mongo{
		User:     "kingangel",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     "27017",
		DB:       "meteor",
		Timeout:  5000,
		Collections: []configs.Collection{
			{
				Name: "node",
			},
		},
	}
	InitMongo(m)
}

func TestNode(t *testing.T) {
	Init()
	nv := models.NodeVO{
		Type:     "SCRIPT",
		Describe: "sdfdsssssss",
	}
	if result, err := CreateNode(nv); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}

func TestUpdateNode(t *testing.T) {
	Init()
	n := models.NodeVO{
		ID:       "61b99e4b73b97742767f9266",
		Type:     "SCRIPT",
		Describe: "rrr33333ss",
		Http: &models.Http{
			Host: "eddd",
		},
	}
	if err := UpdateNode(n); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("update node success")
	}
}

func TestFindNode(t *testing.T) {
	Init()
	id, _ := primitive.ObjectIDFromHex("61b99a439b3a48ffe7a57d7e")
	fmt.Println(GetNode(id))
}

func TestDeleteNode(t *testing.T) {
	Init()
	id, _ := primitive.ObjectIDFromHex("61b99e238a1e886c353f2f34")
	fmt.Println(DeleteNode(id))
}

func Test(t *testing.T) {
	aa()
}

func aa() {
	id, _ := primitive.ObjectIDFromHex("61b99e4b73b97742767f9266")
	n := Node{
		ID:       id,
		Type:     "SCRIPT",
		Describe: "sdfdssssss23423423423sswerewsssssss",
	}
	t := reflect.TypeOf(n)
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Name)
	}
}
