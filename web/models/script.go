package models

import (
	"common"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ScriptNode map[string][]BaseScript

type BaseScript struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Code string `json:"code,omitempty"`
}

type ScriptVO struct {
	BaseScript
	CreateTime string `json:"createTime,omitempty" bson:"createTime,omitempty"`
	UpdateTIme string `json:"updateTime,omitempty" bson:"updateTime,omitempty"`
}

func (sv *ScriptVO) GetBaseBson() bson.M {
	base := bson.M{}
	if sv.Name != "" {
		base["name"] = sv.Name
	}
	if sv.Type != "" {
		base["type"] = sv.Type
	}
	if sv.Code != "" {
		base["code"] = sv.Code
	}
	return base
}

func (sv *ScriptVO) GetUpdateBson() bson.M {
	update := sv.GetBaseBson()
	update["updateTime"] = time.Now().In(common.GetLocalZone())
	return update
}

func (sv *ScriptVO) GetCreateBson() bson.M {
	create := sv.GetBaseBson()
	create["createTime"] = time.Now().In(common.GetLocalZone())
	create["updateTime"] = time.Now().In(common.GetLocalZone())
	return create
}

type ScriptListVO struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}
