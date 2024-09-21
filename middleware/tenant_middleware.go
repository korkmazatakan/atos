// middleware/tenant_middleware.go

package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type TenantConfig struct {
	TenantID string         `json:"tenantID"`
	Database DatabaseConfig `json:"database"`
}

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// origin := c.GetHeader("Origin")
		// u, err := url.Parse(origin)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse origin url"})
		// 	c.Abort()
		// 	return
		// }

		// fmt.Println(u.Host)
		tenantID := c.GetHeader("X-Tenant-ID")

		tenantConfig, err := loadTenantConfig(tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load tenant configuration"})
			c.Abort()
			return
		}

		// Set the tenant information and database connection in the context
		c.Set("tenantConfig", tenantConfig)
		c.Set("db", getDBConnection(tenantConfig.Database))

		c.Next()
	}
}

func loadTenantConfig(tenantID string) (*TenantConfig, error) {
	configPath := filepath.Join("config", "tenants", fmt.Sprintf("%s.json", tenantID))
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var tenantConfig TenantConfig
	err = json.Unmarshal(file, &tenantConfig)
	if err != nil {
		return nil, err
	}

	return &tenantConfig, nil
}

func getDBConnection(dbConfig DatabaseConfig) *gorm.DB {
	dialector := postgres.New(postgres.Config{
		DSN: fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
			dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Host, dbConfig.Port),
	})

	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	return db
}
