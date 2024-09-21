package middleware

import (
	"atos/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DBMigration() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("migration started!")
		defer fmt.Println("migration finished!")
		db, exists := c.Get("db")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Db not exist in context for migration"})
			c.Abort()
			return
		}

		db.(*gorm.DB).AutoMigrate(&models.User{})

		c.Next()
	}
}
