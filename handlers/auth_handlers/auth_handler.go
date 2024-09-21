package handlers

import (
	"atos/models"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginHandler handles user login
func LoginHandler(c *gin.Context) {
    // TODO :: check session if session was set don't override

	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenant information"})
		c.Abort()
		return
	}
	username := c.PostForm("username")
	password := c.PostForm("password")

   //  TODO :: validation check for parameters provided by user
   // username, err := validation.RealEscapeString(username)
   // if err != nil {
   // 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenant information"})
   // 	c.Abort()
   // 	return
   // }

	user := models.User{Username: username}
	res := db.(*gorm.DB).Where("username = ?", username).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) || res.RowsAffected == 0 || res.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		c.Abort()
		return
	}
	passwordCheck := user.CheckPassword(password)
	if !passwordCheck {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		c.Abort()
		return
	}

	session := sessions.Default(c)
	session.Set("authenticated", true)
	session.Set("username", user.Username)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)

	// Retrieve the username from the session
	username := session.Get("username")

	// Clear the session
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful", "username": username})
}

func Register(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenant information"})
		c.Abort()
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	utype := c.PostForm("type")

	if username == "" || password == "" || email == "" || utype == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "send name, email and type as a form data"})
		c.Abort()
		return
	}

	user := models.User{
		Username:  username,
		Email:     email,
		Type:      utype,
		Confirmed: false,
	}
	user.SetPassword(password)
	db.(*gorm.DB).Create(&user)
	c.JSON(http.StatusOK, gin.H{"users": user.ID})
}
