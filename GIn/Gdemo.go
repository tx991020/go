package main

import (
	"github.com/appleboy/gin-jwt"
	"github.com/astaxie/beego/orm"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"

	"fmt"
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

func main() {
	r := gin.Default()
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
				return userId, true
			}

			return userId, true
		},
	}
	r.POST("/users", PostUser)
	r.Use(authMiddleware.MiddlewareFunc())
	v1 := r.Group("api/v1")
	{

		v1.GET("/users", GetUsers)

		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
	}

	r.Run(":8080")
}

func PostUser(c *gin.Context) {
	o := orm.NewOrm()
	var users Users
	users.Firstname = c.PostForm("firstname")
	users.Lastname = c.PostForm("lastname")
	_, err := o.Insert(&users)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Fields are empty"})
	}
	c.JSON(http.StatusOK, gin.H{"result": users})

}

func GetUsers(c *gin.Context) {
	o := orm.NewOrm()
	var (
		users []Users
	)
	_, err := o.QueryTable("users").OrderBy("id").All(&users)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Fields are empty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": users})

}

func GetUser(c *gin.Context) {
	o := orm.NewOrm()
	var users Users
	id := c.Param("id")
	user_id, _ := strconv.Atoi(id)

	err := o.QueryTable("users").Filter("id", user_id).One(&users)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "User not found"})
		return

	}
	c.JSON(http.StatusOK, gin.H{"result": users})

}

func UpdateUser(c *gin.Context) {
	o := orm.NewOrm()
	var users Users
	id := c.Param("id")
	user_id, _ := strconv.Atoi(id)

	users.Firstname = c.PostForm("firstname")
	users.Lastname = c.PostForm("lastname")
	num, err := o.QueryTable("users").Filter("id", user_id).Update(orm.Params{
		"firstname": users.Firstname,
		"lastname":  users.Lastname,
	})
	fmt.Println(num)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "User not found"})
		return

	}

	c.JSON(http.StatusOK, gin.H{"result": users})

}

func DeleteUser(c *gin.Context) {

	o := orm.NewOrm()

	id := c.Param("id")
	user_id, _ := strconv.Atoi(id)
	_, err := o.QueryTable("users").Filter("id", user_id).Delete()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "User not found"})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully deleted user: %s", user_id)})

}
