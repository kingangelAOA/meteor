package routers

import (
	"github.com/gin-gonic/gin"
	"meteor/controllers"
)

//InitRoutes init routes
func InitRoutes(route *gin.RouterGroup) {
	login := route.Group("/user")
	login.POST("/login", controllers.Login)
	login.GET("/info", controllers.Info)

	project := route.Group("/project")
	project.POST("", controllers.CreateOrUpdateProject)
	project.GET("/all", controllers.GetAllProjects)

	swagger := route.Group("/swagger")
	swagger.GET("/versions", controllers.GetSwaggerVersions)
	swagger.POST("", controllers.UpdateSwagger)
	swagger.GET("", controllers.GetSwaggerYaml)
}
