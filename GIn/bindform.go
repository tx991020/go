package main

import (
	"github.com/gin-gonic/gin"
	"net/http"

)

type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.POST("/loginJSON", func(c *gin.Context) {

		var json Login
		if c.Bind(&json) == nil {

			if json.User == "manu" && json.Password == "123" {
				// INSERT INTO "users" (name) VALUES (user.Name);

				// Display error
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				// Display error
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})

	r.Run(":8080")
}
