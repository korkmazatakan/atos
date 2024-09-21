package main

import (
	authhandler "atos/handlers/auth_handlers"
	userhandlers "atos/handlers/user_handlers"
	"atos/middleware"
	"net/http"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.TenantMiddleware())
	router.Use(middleware.SessionMiddleware())
	router.Use(middleware.DBMigration())

	router.POST("/register", authhandler.Register)
	router.POST("/login", authhandler.LoginHandler)
	router.GET("/logout", authhandler.LogoutHandler)

	authGroup := router.Group("/auth")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/dashboard", func(c *gin.Context) {
			username := sessions.Default(c).Get("username")
			c.JSON(http.StatusOK, gin.H{"message": "Welcome " + username.(string)})
		})
	}

	router.GET("/users", userhandlers.GetUsers)
	router.GET("/users/me", userhandlers.GetMe)
	router.GET("/users/:id", userhandlers.GetById)
	router.POST("/users/add", userhandlers.AddUser)

	router.Run(":8080")
}
