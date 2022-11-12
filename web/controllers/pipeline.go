package controllers

import (
	"common/bindmodels"
	"web/models"
	"web/service"

	"github.com/gin-gonic/gin"
)

func GetPipelines(c *gin.Context) {
	var pq bindmodels.Query
	if err := c.BindQuery(&pq); err != nil {
		if err != nil {
			ResponseErr(c, err.Error())
			return
		}
	}
	vos, err := service.GetPipelines(pq)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, vos)
}

func GetPipelinePagination(c *gin.Context) {
	var pq bindmodels.Query
	if err := c.BindQuery(&pq); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	vos, err := service.GetPipelinePagination(pq)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, vos)
}

func CreatePipeline(c *gin.Context) {
	var pv models.PipelineVO
	if err := c.BindJSON(&pv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	if _, err := service.CreatePipeline(&pv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func UpdatePipeline(c *gin.Context) {
	var pv models.PipelineVO
	if err := c.BindJSON(&pv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	if err := service.UpdatePipeline(&pv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func DeletePipeline(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		if err := service.DeletePipeline(id); err != nil {
			ResponseErr(c, err.Error())
			return
		}
		ResponseSuccessNoData(c)
		return
	}
	ResponseErr(c, "id is not exist")
}

func UpdateFlow(c *gin.Context) {
	var fv models.FlowVO
	if err := c.BindJSON(&fv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	if err := service.UpdateFlow(&fv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func GetFlowByPipelineId(c *gin.Context) {
	if pipelineId, ok := c.GetQuery("pipelineId"); ok {
		if f, err := service.GetFlowByPipelineId(pipelineId); err != nil {
			ResponseErr(c, err.Error())
			return
		} else {
			ResponseSuccess(c, f)
			return
		}
	}
	ResponseErr(c, "pipelineId is not exist")
}

func RunPipeline(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		taskId, err := service.RunPipeline(id)
		if err != nil {
			ResponseErr(c, err.Error())
			return
		}
		ResponseSuccess(c, taskId)
		return
	}
	ResponseErr(c, "id is not exist")
}

func GetPipelineTypes(c *gin.Context)  {
	ResponseSuccess(c, service.GetPipelineTypes())
}