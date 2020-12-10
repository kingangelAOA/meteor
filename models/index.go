package models

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//Response struct
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//ResponseErr 错误返回
func ResponseErr(c *gin.Context, err string) {
	c.JSON(http.StatusOK, Response{
		Code: 500,
		Msg:  err,
	})
}

//ResponseSuccess 正确带数据返回
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: data,
	})
}

//ResponseSuccessNoData 正确不带数据返回
func ResponseSuccessNoData(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "ok",
	})
}
