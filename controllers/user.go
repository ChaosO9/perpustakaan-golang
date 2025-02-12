package controllers

import (
	"perpustakaan-golang/forms"
	"perpustakaan-golang/models"
	"perpustakaan-golang/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

//UserController ...
type UserController struct{}

var userModel = new(models.UserModel)
var userForm = new(forms.UserForm)

//getUserID ...
func getUserID(c *gin.Context) (userID int64) {
	//MustGet returns the value for the given key if it exists, otherwise it panics.
	return c.MustGet("userID").(int64)
}

//Login ...
// Login handles user login.
// @Summary User Login
// @Description Logs in a user.
// @Tags User
// @Accept json
// @Produce json
// @Param login body forms.LoginForm true "Login Request"
// @Success 200 {object} models.LoginResponse  
// @Failure 406 {object} models.ErrorResponse
// @Router /user/login [post]
func (ctrl UserController) Login(c *gin.Context) {
	var loginForm forms.LoginForm

	if validationErr := c.ShouldBindJSON(&loginForm); validationErr != nil {
		message := userForm.Login(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	user, token, err := userModel.Login(loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid login details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "user": user, "token": token})
}

//Register ...
// Register handles user registration.
// @Summary User Register
// @Description Register a new user
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param register body forms.RegisterFormSwagger true "Register Request"
// @Success 200 {object} models.RegisterResponse
// @Failure 406 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/register [post]
func (ctrl UserController) Register(c *gin.Context) {
	var registerForm forms.RegisterForm

	if validationErr := c.ShouldBind(&registerForm); validationErr != nil {
		message := userForm.Register(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	uploadDir := "./uploads/user/"
    registerForm.FotoFileName = "profile-placeholder.jpg"

    filename, err := utils.ImageUpload(c, "foto", uploadDir) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	registerForm.FotoFileName = filename

	user, err := userModel.Register(registerForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered", "user": user})
}

//Logout ...
// Logout handles user logout.
// @Summary User Logout
// @Description Logs out a user.
// @Tags User
// @Security ApiKeyAuth  
// @Success 200 {object} models.LogoutResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /user/logout [get]
func (ctrl UserController) Logout(c *gin.Context) {

	au, err := authModel.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not logged in"})
		return
	}

	deleted, delErr := authModel.DeleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}