package routers

import (
	"web/controllers"

	"github.com/gin-gonic/gin"
)

// InitRoutes init routes
func InitRoutes(route *gin.RouterGroup) {
	login := route.Group("/user")
	login.POST("/login", controllers.Login)
	login.GET("/info", controllers.Info)

	// project := route.Group(	"/project")
	// project.POST("", controllers.CreateOrUpdateProject)
	// project.GET("/all", controllers.GetAllProjects)
	pipeline := route.Group("/pipeline")
	pipeline.POST("", controllers.CreatePipeline)
	pipeline.PUT("", controllers.UpdatePipeline)
	pipeline.DELETE("", controllers.DeletePipeline)
	pipeline.GET("/list", controllers.GetPipelinePagination)
	pipeline.PUT("/flow", controllers.UpdateFlow)
	pipeline.GET("/flow", controllers.GetFlowByPipelineId)
	pipeline.POST("/run", controllers.RunPipeline)
	pipeline.GET("/tasks", controllers.GetTaskPagination)
	pipeline.GET("/types", controllers.GetPipelineTypes)

	task := route.Group("/task")
	task.GET("/running/flow/info", controllers.GetTaskRunningFlowInfo)
	task.GET("/running/monitor/info", controllers.GetTaskRunningMonitorInfo)
	task.POST("/stop", controllers.StopTask)
	task.POST("/reset/qps", controllers.ResetTaskQPS)
	task.POST("/reset/users", controllers.ResetTaskUsers)
	task.GET("/pipeline/users", controllers.GetTaskPipelineUsers)
	task.GET("/status", controllers.GetTaskStatus)
	task.GET("/works", controllers.GetRunningWorks)

	node := route.Group("/node")
	node.GET("/list", controllers.GetNodePagination)
	node.GET("/all", controllers.GetNodes)
	node.GET("/binding/plugin", controllers.GetBindingPluginNodes)
	node.POST("", controllers.CreateNode)
	node.PUT("", controllers.UpdateNode)
	node.DELETE("", controllers.DeleteNode)

	plugin := route.Group("/plugin")
	plugin.GET("/list", controllers.GetPluginPagination)
	plugin.GET("/nodes/all", controllers.GetPlugins)
	plugin.POST("", controllers.CreatePlugin)
	plugin.PUT("", controllers.UpdatePlugin)
	plugin.DELETE("", controllers.DeletePlugin)
	plugin.GET("", controllers.GetPluginByID)
	plugin.PUT("/file", controllers.UpdatePluginFile)
	plugin.POST("/debug", controllers.DebugPlugin)
	plugin.GET("/selectList", controllers.GetPluginSelectList)
	plugin.GET("/file/status", controllers.CheckPluginFileStatus)
	plugin.POST("/publish", controllers.PublishPlugin)

	script := route.Group("/script")
	script.GET("/list", controllers.GetScriptList)

	// component := route.Group("component")
	// component.GET("components", controllers.GetComponents)
	// component.GET("component", controllers.GetComponentByID)
	// component.POST("component", controllers.CreateComponent)
	// component.PUT("component", controllers.UpdateComponent)
	// component.DELETE("component", controllers.DeleteComponent)
}
