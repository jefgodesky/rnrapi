package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/models"
)

func IsWorldCreator(world *models.World, user *models.User) bool {
	if user == nil || world == nil {
		return false
	}

	for _, creator := range world.Creators {
		if creator.ID == user.ID {
			return true
		}
	}

	return false
}

func HasWorldAccess(world *models.World, user *models.User) bool {
	if world.Public {
		return true
	}

	if IsWorldCreator(world, user) {
		return true
	}

	return false
}

func WorldCreatorOnly(c *gin.Context) *models.World {
	world := GetWorldFromSlug(c)
	if world == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	isCreator := IsWorldCreator(world, user)
	if !isCreator {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return world
}

func FilterCampaignWorldAccess(campaigns []models.Campaign, user *models.User) []models.Campaign {
	var filtered []models.Campaign
	for _, campaign := range campaigns {
		if HasWorldAccess(&campaign.World, user) {
			filtered = append(filtered, campaign)
		}
	}
	return filtered
}

func IsCampaignGM(campaign *models.Campaign, user *models.User) bool {
	if user == nil || campaign == nil {
		return false
	}

	for _, gm := range campaign.GMs {
		if gm.ID == user.ID {
			return true
		}
	}

	return false
}

func HasCampaignAccess(campaign *models.Campaign, user *models.User) bool {
	if !HasWorldAccess(&campaign.World, user) {
		return false
	}

	if campaign.Public {
		return true
	}

	if IsCampaignGM(campaign, user) {
		return true
	}

	return false
}

func CampaignGMOnly(c *gin.Context) *models.Campaign {
	campaign := GetCampaignFromSlug(c)
	if campaign == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	if !HasCampaignAccess(campaign, user) || !IsCampaignGM(campaign, user) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return campaign
}
