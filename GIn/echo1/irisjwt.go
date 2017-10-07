package main

import (
	"github.com/kataras/iris"

	jwt "github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"time"
)

type User struct {
	Id        int
	Password string
	Username  string
}


func myHandler(ctx iris.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)


	ctx.Writef("%s", "Welcome "+name+"!")
}

func main() {
	app := iris.New()

	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})
	app.Post("/login", login)
	app.Use(jwtHandler.Serve)

	app.Get("/restricted", myHandler)
	app.Run(iris.Addr("localhost:3001"))
	
}

func login(ctx iris.Context)  {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	var user User
	if username == user.Username && password == user.Password{

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"name":"Jon Snow",
			"admin":true,
			"timeout":jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		})
		t, _ := token.SignedString([]byte("My Secret"))

		ctx.JSON(iris.Map{"token": t})
	}else {
		//ctx.WriteString(" username or password  wrong")
		ctx.Redirect("http://127.0.0.1:3001")
	}
}