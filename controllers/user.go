package controllers

import (
	"meteor/models"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	models.ResponseSuccess(c, models.UserRes{
		Token: "admin-token",
	})
}

func Info(c *gin.Context) {
	models.ResponseSuccess(c, models.UserInfoRes{
		Roles:        []string{"admin"},
		Introduction: "am a super administrator",
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Name:         "Super Admin",
	})
}
