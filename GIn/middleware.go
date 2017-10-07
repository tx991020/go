package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"net/http"
	"fmt"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 在gin上下文中定义变量
		c.Set("example", "12345")

		// 请求前

		c.Next() //处理请求

		// 请求后
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func main()  {
	r := gin.New()
	r.Use(Logger())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "It works")
	})
	r.POST("/post", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})
	r.GET("/test", func(c *gin.Context) {
		//获取gin上下文中的变量
		example := c.MustGet("example").(string)

		// 会打印: "12345"
		c.String(http.StatusOK,example)

	})

	r.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	// 监听运行于 0.0.0.0:8080
	r.Run(":8080")

}