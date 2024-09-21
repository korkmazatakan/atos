package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins: []string{"*"},
		// AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders: []string{"*"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin != "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	})
}
