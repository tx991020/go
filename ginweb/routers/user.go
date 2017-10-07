package routers

import (
	"github.com/gin-gonic/gin"
	"taskmanager/controllers"
)

func SetUserRoutes(router *gin.Engine) *gin.Engine {

	userR := router.Group("/tm2/users")
	{
		userR.POST("/register", controllers.Register)
		userR.POST("/login", controllers.Login)
	}
	return router
}
