package controllers

import (
	"common"
	"common/bindmodels"
	"web/models"
	"web/service"

	"github.com/gin-gonic/gin"
)

func GetNodes(c *gin.Context) {
	var q bindmodels.Query
	err := c.BindQuery(&q)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	vos, err := service.GetNodes(q)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, vos)
}

func GetBindingPluginNodes(c *gin.Context) {
	var q bindmodels.Query
	err := c.BindQuery(&q)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	vos, err := service.GetBindingPluginNodes(q)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, vos)
}

func GetAllNodes(c *gin.Context) {
	vos, err := service.GetAllNodes()
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, vos)
}

func GetNodePagination(c *gin.Context) {
	var pq bindmodels.Query
	if err := c.BindQuery(&pq); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	vos, err := service.GetNodePagination(pq)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, vos)
}

func CreateNode(c *gin.Context) {
	var nv models.NodeVO
	if err := c.BindJSON(&nv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	nv.Type = common.Plugin
	if _, err := service.CreateNode(&nv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func UpdateNode(c *gin.Context) {
	var nv models.NodeVO
	if err := c.BindJSON(&nv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	if err := service.UpdateNode(&nv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func DeleteNode(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		if err := service.DeleteNode(id); err != nil {
			ResponseErr(c, err.Error())
			return
		}
		ResponseSuccessNoData(c)
		return
	}
	ResponseErr(c, "id is not exist")
}
