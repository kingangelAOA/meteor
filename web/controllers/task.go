package controllers

import (
	"common/bindmodels"
	"web/service"

	"github.com/gin-gonic/gin"
)

func GetTaskPagination(c *gin.Context) {
	var pq bindmodels.Query
	if err := c.BindQuery(&pq); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ts, err := service.GetTaskPagination(pq)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, ts)
}

func GetTaskRunningFlowInfo(c *gin.Context) {
	var tq bindmodels.TaskQuery
	if err := c.BindQuery(&tq); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	tv, err := service.GetTaskRunningFlowInfo(tq)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	tv.Status = service.GetTaskStatus(tv.ID)
	ResponseSuccess(c, tv)
}

func GetTaskRunningMonitorInfo(c *gin.Context) {
	var tq bindmodels.TaskQuery
	if err := c.BindQuery(&tq); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	trm, err := service.GetTaskRunningMonitorInfo(tq)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, trm)
}

func GetTaskStatus(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		ResponseSuccess(c, service.GetTaskStatus(id))
		return
	}
	ResponseErr(c, "id is not exist")
}

func GetRunningWorks(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		works, err := service.GetRunningWorks(id)
		if err != nil {
			ResponseErr(c, err.Error())
			return
		}
		ResponseSuccess(c, works)
		return
	}
	ResponseErr(c, "id is not exist")
}

func StopTask(c *gin.Context) {
	var too bindmodels.TaskOperateOption
	if err := c.BindJSON(&too); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	err := service.StopTask(too.ID)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func ResetTaskQPS(c *gin.Context) {
	var too bindmodels.TaskOperateOption
	if err := c.BindJSON(&too); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	err := service.ResetQPS(too.ID, too.QPS)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func ResetTaskUsers(c *gin.Context) {
	var too bindmodels.TaskOperateOption
	if err := c.BindJSON(&too); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	err := service.ResetUsers(too.ID, too.Users)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func GetTaskPipelineUsers(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		users, err := service.GetTaskPipelineUsers(id)
		if err != nil {
			ResponseErr(c, err.Error())
			return
		}
		ResponseSuccess(c, users)
		return
	}
	ResponseErr(c, "id is not exist")
}
