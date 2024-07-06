package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
	"gorm.io/gorm/clause"
)

func CampaignCreate(c *gin.Context) {
	campaign := helpers.BodyToCampaign(c)

	if result := initializers.DB.Create(&campaign); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create campaign"})
		return
	}

	c.JSON(200, serializers.SerializeCampaign(*campaign))
}

func CampaignIndex(c *gin.Context) {
	var campaigns []models.Campaign
	user := helpers.GetUserFromContext(c, false)

	if user != nil {
		initializers.DB.
			Preload(clause.Associations).
			Where("public = ? OR id in (SELECT campaign_id FROM campaign_gms WHERE user_id = ?)", true, user.ID).
			Find(&campaigns)
	} else {
		initializers.DB.
			Preload(clause.Associations).
			Where("Public = ?", true).
			Find(&campaigns)
	}

	filtered := helpers.FilterCampaignWorldAccess(campaigns, user)
	c.JSON(200, gin.H{
		"campaigns": serializers.SerializeCampaigns(filtered),
	})
}
