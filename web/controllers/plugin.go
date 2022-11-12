package controllers

import (
	"common/bindmodels"
	"web/models"
	"web/service"

	"github.com/gin-gonic/gin"
)

func GetPlugins(c *gin.Context) {
	var q bindmodels.Query
	c.BindQuery(&q)
	vos, err := service.GetPluginPagination(q)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, vos)
}

func GetPluginPagination(c *gin.Context) {
	var pq bindmodels.Query
	if err := c.BindQuery(&pq); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	vos, err := service.GetPluginPagination(pq)
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, vos)
}

func CreatePlugin(c *gin.Context) {
	var pv models.PluginVO
	if err := c.BindJSON(&pv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	if oid, err := service.CreatePlugin(&pv); err != nil {
		ResponseErr(c, err.Error())
		return
	} else {
		ResponseSuccess(c, oid.Hex())
	}
}

func GetPluginByID(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		ResponseSuccess(c, models.PluginVO{})
		return
	}
	inputs, err := service.GetPluginByID(c.Query("id"))
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, inputs)
}

func UpdatePlugin(c *gin.Context) {
	var pv models.PluginVO
	if err := c.BindJSON(&pv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	if err := service.UpdatePlugin(&pv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func DeletePlugin(c *gin.Context) {
	if id, ok := c.GetQuery("id"); ok {
		if err := service.DeletePlugin(id); err != nil {
			ResponseErr(c, err.Error())
			return
		}
		ResponseSuccessNoData(c)
		return
	} else {
		ResponseErr(c, "id is not exist")
	}
}

func UpdatePluginFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	src, err := file.Open()
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	service.ReceiveBinaryFile(data, c.PostForm("id"))
    err = src.Close()
    if err != nil {
        ResponseErr(c, err.Error())
        return
    }
	ResponseSuccessNoData(c)
}

func DebugPlugin(c *gin.Context) {
	var pv models.PluginVO
	if err := c.BindJSON(&pv); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	if err := service.UpdatePlugin(&pv); err != nil {
		ResponseSuccess(c, err.Error())
		return
	}
	ResponseSuccess(c, service.DebugPlugin(pv))
}

func GetPluginSelectList(c *gin.Context) {
	list, err := service.GetPluginList()
	if err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccess(c, list)
}

func CheckPluginFileStatus(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		ResponseErr(c, "id is not exist")
		return
	}
	language, ok := c.GetQuery("language")
	if !ok {
		ResponseErr(c, "language is not exist")
		return
	}
	if err := service.CheckPluginFileStatus(id, language); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}

func PublishPlugin(c *gin.Context) {
	var pp bindmodels.PluginParams
	if err := c.BindJSON(&pp); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	if err := service.PublishPlugin(pp.ID); err != nil {
		ResponseErr(c, err.Error())
		return
	}
	ResponseSuccessNoData(c)
}
