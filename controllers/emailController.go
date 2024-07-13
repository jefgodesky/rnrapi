package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/parsers"
	"github.com/jefgodesky/rnrapi/serializers"
)

func EmailCreate(c *gin.Context) {
	email := parsers.BodyToEmail(c)
	if email == nil {
		return
	}

	result := initializers.DB.Create(&email)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Failed to create email"})
		c.Abort()
		return
	}

	c.JSON(200, serializers.SerializeEmail(*email))
}
