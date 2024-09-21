package userhandlers

import (
	"atos/models"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUsers(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenant information"})
		c.Abort()
		return
	}

	users := []models.User{}
	db.(*gorm.DB).Find(&users)

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetMe(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenant information"})
		c.Abort()
		return
	}
	session := sessions.Default(c)
	username := fmt.Sprintf("%v", session.Get("username"))

	user := models.User{}
	res := db.(*gorm.DB).Where("username = ?", username).First(&user)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user couldn't be found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetById(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenant information"})
		c.Abort()
		return
	}

	id := c.Param("id")

	user := models.User{}
	res := db.(*gorm.DB).First(&user, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user couldn't be found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetByUsername(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenant information"})
		c.Abort()
		return
	}

	username := c.Param("username")

	user := models.User{}
	res := db.(*gorm.DB).Where("username = ?", username).First(&user)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user couldn't be found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func AddUser(c *gin.Context) {
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
