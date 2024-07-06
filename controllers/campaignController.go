package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/serializers"
)

func CampaignCreate(c *gin.Context) {
	campaign := helpers.BodyToCampaign(c)

	if result := initializers.DB.Create(&campaign); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create campaign"})
		return
	}

	c.JSON(200, serializers.SerializeCampaign(*campaign))
}
