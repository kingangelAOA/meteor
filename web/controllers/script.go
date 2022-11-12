package controllers

import (
	"web/service"

	"github.com/gin-gonic/gin"
)

func GetScriptList(c *gin.Context) {
	ls := service.GetScriptList()
	ResponseSuccess(c, ls)
}
