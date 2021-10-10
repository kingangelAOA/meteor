package db

import (
	"fmt"
	"meteor/core"
	"testing"
)

func TestStage(t *testing.T) {
	Init()
	ctx, cancel, s := getCollection("stage")
	defer cancel()
	result, err := s.InsertOne(ctx, &Stage{Type: core.MultipleGoroutine})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result.InsertedID)
	}
	ctx, cancel, n := getCollection("node")
	defer cancel()
	id, err := n.InsertOne(ctx, &Node{Name: "test", Type: "SCRIPT"})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(id)
	}
}
