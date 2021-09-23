package middlewares

import (
	"meteor/common"
	"meteor/core"
	. "meteor/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GlobalRecover(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			if XHR(c) {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": rec,
				})
			} else {
				ResponseErr(c, common.FatalError.Error())
			}
		}
	}(c)
	//c.Next()
}

func GlobalResponse(c *gin.Context) {
	if len(c.Errors) > 0 {
		lastErr := c.Errors[len(c.Errors)-1]
		core.ResponseErr(c, lastErr, common.GetCodeByErr(lastErr))
	} else {
		if d, ok := c.Keys[common.ResponseDataKey]; ok {
			core.ResponseSuccess(c, d)
		} else {
			ResponseSuccessNoData(c)
		}
	}
}

func XHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
