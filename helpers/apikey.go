package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/models"
)

func GenerateAPIKey(c *gin.Context) (string, string, string) {
	token, secret, err := models.GenerateAPIKey()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to generate API key"})
		c.Abort()
		return "", "", ""
	}

	hash, err := Hash(secret)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to hash API key"})
		c.Abort()
		return "", "", ""
	}

	return token, hash, token + "." + secret
}
