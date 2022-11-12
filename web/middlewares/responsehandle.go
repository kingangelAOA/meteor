package middlewares

import (
	"net/http"
	"strings"
	"web/common"
	"web/controllers"

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
				controllers.ResponseErr(c, common.FatalError.Error())
			}
		}
	}(c)
	c.Next()
}

func GlobalResponse(c *gin.Context) {
	if len(c.Errors) > 0 {
		lastErr := c.Errors[len(c.Errors)-1]
		lastErr.Err.Error()
		controllers.ResponseErr(c, lastErr.Err.Error())
	}
}

func XHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
