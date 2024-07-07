package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/serializers"
)

func CharacterCreate(c *gin.Context) {
	character := helpers.BodyToCharacter(c)
	if character == nil {
		return
	}

	if result := initializers.DB.Create(&character); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create character"})
		return
	}

	c.JSON(200, serializers.SerializeCharacter(*character))
}
