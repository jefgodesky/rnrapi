package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/serializers"
)

func ScrollCreate(c *gin.Context) {
	scroll := helpers.BodyToScroll(c)

	if result := initializers.DB.Create(&scroll); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create scroll"})
		return
	}

	c.JSON(200, serializers.SerializeScroll(*scroll))
}
