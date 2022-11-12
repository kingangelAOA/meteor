package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func main() {
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			Code: 0,
			Data: map[string]string{
				"ok": "xxxx",
			},
		})
	})
	r.Run("0.0.0.0:7090")
}
