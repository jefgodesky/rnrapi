package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"strings"
)

func getUserFromAPIKey(c *gin.Context, required bool) {
	apiKey := c.GetHeader("Authorization")
	if apiKey == "" || !strings.HasPrefix(apiKey, "Bearer ") {
		if required == true {
			c.JSON(401, gin.H{"error": "API key required"})
			c.Abort()
		}
		c.Next()
		return
	}
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")

	parts := strings.Split(apiKey, ".")
	if len(parts) != 2 {
		if required == true {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
		}
		c.Next()
		return
	}
	token, secret := parts[0], parts[1]

	var key models.Key
	result := initializers.DB.
		Preload("User").
		Where("token = ?", token).First(&key)
	if result.Error != nil {
		if required {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		c.Next()
		return
	}

	if err := helpers.CheckHash(secret, key.Secret); err != nil {
		if required == true {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
		}
		c.Next()
		return
	}

	c.Set("keyID", key.ID)
	c.Set("user", &key.User)
	c.Next()
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		getUserFromAPIKey(c, true)
	}
}

func AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		getUserFromAPIKey(c, false)
	}
}
