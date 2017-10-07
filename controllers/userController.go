package controllers

import (

	"net/http"
	"github.com/gin-gonic/gin"
	"taskmanager/common"
	"taskmanager/data"
	"taskmanager/models"
)

// Register add a new User document
// Handler for HTTP Post - "/users/register"
func Register(c *gin.Context) {
	var dataResource UserResource
	// Decode the incoming User json
	err := c.BindJSON(&dataResource)
	if err != nil {
		common.DisplayAppError(
			c,
			err,
			"Invalid User data",
			500,
		)
		return
	}
	user := &dataResource.Data
	context := NewContext()
	defer context.Close()
	col := context.DbCollection("users")
	repo := &data.UserRepository{C: col}
	// Insert User document
	repo.CreateUser(user)
	// Clean-up the hashpassword to eliminate it from response JSON
	user.HashPassword = nil
	c.JSON(http.StatusCreated,dataResource)

}

// Login authenticates the HTTP request with username and apssword
// Handler for HTTP Post - "/users/login"
func Login(c *gin.Context) {
	var dataResource LoginResource
	var token string
	// Decode the incoming Login json
	err := c.BindJSON(&dataResource)
	if err != nil {
		common.DisplayAppError(
			c,
			err,
			"Invalid Login data",
			500,
		)
		return
	}
	loginModel := dataResource.Data
	loginUser := models.User{
		Email:    loginModel.Email,
		Password: loginModel.Password,
	}
	context := NewContext()
	defer context.Close()
	col := context.DbCollection("users")
	repo := &data.UserRepository{C: col}
	// Authenticate the login user
	user, err := repo.Login(loginUser)
	if err != nil {
		common.DisplayAppError(
			c,
			err,
			"Invalid login credentials",
			401,
		)
		return
	}
	// Generate JWT token
	token, err = common.GenerateJWT(user.Email, "member")
	if err != nil {
		common.DisplayAppError(
			c,
			err,
			"Eror while generating the access token",
			500,
		)
		return
	}

	// Clean-up the hashpassword to eliminate it from response JSON
	user.HashPassword = nil
	authUser := AuthUserModel{
		User:  user,
		Token: token,
	}
	c.JSON(http.StatusOK, AuthUserResource{Data: authUser})
}
