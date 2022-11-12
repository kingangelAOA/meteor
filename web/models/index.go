package models

import (
	"common"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type TypeVO interface {
	PipelineVO | PluginVO
}

func GetUpdateBson(update bson.M) bson.M {
	update["updateTime"] = time.Now().In(common.GetLocalZone())
	return update
}

func GetCreateBson(create bson.M) bson.M {
	create["createTime"] = time.Now().In(common.GetLocalZone())
	create["updateTime"] = time.Now().In(common.GetLocalZone())
	return create
}
