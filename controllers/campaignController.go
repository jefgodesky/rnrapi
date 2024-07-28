package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/parsers"
	"github.com/jefgodesky/rnrapi/serializers"
	"gorm.io/gorm/clause"
)

func CampaignCreate(c *gin.Context) {
	campaign := parsers.BodyToCampaign(c)

	if result := initializers.DB.Create(&campaign); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create campaign"})
		return
	}

	c.JSON(200, serializers.SerializeCampaign(*campaign))
}

func CampaignIndex(c *gin.Context) {
	var campaigns []models.Campaign
	user := helpers.GetUserFromContext(c, false)
	worldSlug := c.Query("world")
	query := initializers.DB.
		Model(&models.Campaign{}).
		Preload(clause.Associations).
		Joins("JOIN worlds ON worlds.id = campaigns.world_id")

	if worldSlug != "" {
		query = query.Where("worlds.slug = ?", worldSlug)
	}

	if user != nil {
		query = query.Where("(campaigns.public = ? AND worlds.public = ?) OR campaigns.world_id IN (SELECT world_id FROM world_creators WHERE user_id = ?)", true, true, user.ID)
	} else {
		query = query.Where("campaigns.public = ? AND worlds.public = ?", true, true)
	}

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&campaigns)

	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"campaigns": serializers.SerializeCampaigns(campaigns),
	})
}

func CampaignRetrieve(c *gin.Context) {
	campaign := helpers.GetCampaignFromSlug(c)
	user := helpers.GetUserFromContext(c, false)
	allowed := helpers.HasCampaignAccess(campaign, user)

	if !allowed {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeCampaign(*campaign))
}

func CampaignUpdate(c *gin.Context) {
	campaign := helpers.CampaignGMOnly(c)
	if campaign == nil {
		return
	}

	newCampaign := parsers.BodyToCampaign(c)
	campaign.Slug = newCampaign.Slug
	campaign.Name = newCampaign.Name
	campaign.Description = newCampaign.Description
	campaign.GMs = newCampaign.GMs
	campaign.Public = newCampaign.Public
	campaign.WorldID = newCampaign.WorldID
	campaign.World = newCampaign.World

	if err := initializers.DB.Save(campaign).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update campaign"})
		return
	}

	c.JSON(200, serializers.SerializeCampaign(*campaign))
}

func CampaignDestroy(c *gin.Context) {
	campaign := helpers.CampaignGMOnly(c)
	if campaign == nil {
		return
	}

	if err := initializers.DB.Delete(&campaign).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to destroy campaign"})
		return
	}

	c.JSON(200, serializers.SerializeCampaign(*campaign))
}
