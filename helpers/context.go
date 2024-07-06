package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/models"
)

func GetUserFromContext(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(500, gin.H{"error": "Failed to retrieve user from context"})
		c.Abort()
		return nil
	}
	return user.(*models.User)
}
