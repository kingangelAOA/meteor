package db

import (
	"fmt"
	"meteor/common"
	"meteor/models"
	"testing"
)

func TestPipeline(t *testing.T) {
	Init()
	p := models.PipelineVO{}
	p.Describe = "xxxxsdfsdfx"
	p.Type = "MultipleGoroutine"
	if result, err := CreatePipeline(p); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}

func TestAddNode(t *testing.T) {
	Init()
	oid, _ := common.ConverObjectID("61bc5eaf5e5827c04cf10383")
	nid, _ := common.ConverObjectID("61b964a1b5b260a1284824f5")
	if err := AddNode(oid, nid); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("add node success")
}

func TestDeleteNodeFromPipeline(t *testing.T) {
	Init()
	oid, _ := common.ConverObjectID("61b999c024ba039cb40a4cb3")
	nid, _ := common.ConverObjectID("61b99a439b3a48ffe7a57d7e")
	if err := DeleteNodeFromPipeline(oid, nid); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("delete node success")
}

func TestGetPipelines(t *testing.T) {
	Init()
	r, _ := GetPipelines(10, 0, "createTime", -1)
	fmt.Println(r)
}
