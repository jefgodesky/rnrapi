package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/serializers"
)

func SocietyCreate(c *gin.Context) {
	society := helpers.BodyToSociety(c)
	if society == nil {
		return
	}

	if result := initializers.DB.Create(&society); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create society"})
		return
	}

	c.JSON(200, serializers.SerializeSociety(*society))
}
