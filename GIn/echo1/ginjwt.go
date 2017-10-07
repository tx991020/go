package main

import (
	"fmt"
	"github.com/appleboy/gin-jwt"
	"github.com/astaxie/beego/orm"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Users struct {
	Id        int
	Firstname string
	Lastname  string
}

func init() {
	//注册驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(new(Users))

	orm.RegisterDataBase("default", "mysql", "root:111111@tcp(127.0.1:3306)/test?charset=utf8")
	//自动创建表
	orm.RunSyncdb("default", false, true)
}

func haoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

func main() {
	r := gin.New()
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			var users Users
			o := orm.NewOrm()
			err := o.QueryTable("users").Filter("firstname", userId).One(&users)
			if err != nil {
				fmt.Print(err)
				return "not found", false
			}
			fmt.Println(users)
			if userId == users.Firstname && password == users.Lastname {
				return userId, true
			}

			return userId, false
		},
	}
	r.POST("/login", authMiddleware.LoginHandler)
	r.Use(authMiddleware.MiddlewareFunc())
	r.GET("/hello", haoHandler)
	//r.GET("/refresh_token", authMiddleware.RefreshHandler)
	r.Run(":8000")
}
