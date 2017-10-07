package routers

import (
	"github.com/gin-gonic/gin"
	"taskmanager/common"
	"taskmanager/controllers"
)

func SetTaskRoutes(router *gin.Engine) *gin.Engine {

	taR := router.Group("/tm2/tasks")
	taR.Use(common.Authorize())
	{
		taR.POST("", controllers.CreateTask)
		taR.PUT(":id", controllers.UpdateTask)
		taR.DELETE(":id", controllers.DeleteTask)
		taR.GET("", controllers.GetTasks)
		taR.GET("t/:id", controllers.GetTaskByID)
		taR.GET("users/:email", controllers.GetTasksByUser)
	}
	return router
}
