package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"strings"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
		if apiKey == "" || !strings.HasPrefix(apiKey, "Bearer ") {
			c.JSON(401, gin.H{"error": "API key required"})
			c.Abort()
			return
		}
		apiKey = strings.TrimPrefix(apiKey, "Bearer ")

		parts := strings.Split(apiKey, ".")
		if len(parts) != 2 {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		providedToken, providedSecret := parts[0], parts[1]

		var user models.User
		result := initializers.DB.Where("token = ?", providedToken).First(&user)
		if result.Error != nil {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		if err := models.CheckAPIKey(providedSecret, user.Secret); err != nil {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}
