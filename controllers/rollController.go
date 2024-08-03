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

func RollCreate(c *gin.Context) {
	roll := parsers.BodyToRoll(c)
	if roll == nil {
		return
	}

	helpers.RollOnTable(roll.Table, roll, roll.Modifier, roll.Character)

	if result := initializers.DB.Create(&roll); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create roll record"})
		return
	}

	c.JSON(200, serializers.SerializeRoll(*roll))
}

func RollIndex(c *gin.Context) {
	var rolls []models.Roll
	user := helpers.GetUserFromContext(c, true)
	query := initializers.DB.Model(&models.Roll{})
	query.Joins("LEFT JOIN campaigns ON rolls.campaign_id = campaigns.id").
		Joins("LEFT JOIN campaign_gms ON campaign_gms.campaign_id = campaigns.id").
		Where("rolls.roller_id = ? OR campaign_gms.user_id = ?", user.ID, user.ID).
		Preload(clause.Associations)

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&rolls)

	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"rolls":     serializers.SerializeRolls(rolls),
	})
}

func RollRetrieve(c *gin.Context) {
	roll := helpers.GetRollFromID(c)
	if roll == nil {
		return
	}

	user := helpers.GetUserFromContext(c, false)
	isRoller := roll.Roller.ID == user.ID
	isGM := helpers.IsCampaignGM(roll.Campaign, user)

	if !isRoller && !isGM {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeRoll(*roll))
}

func RollDestroy(c *gin.Context) {
	roll := helpers.GetRollFromID(c)
	if roll == nil {
		return
	}

	user := helpers.GetUserFromContext(c, false)
	isRoller := roll.Roller.ID == user.ID
	isGM := helpers.IsCampaignGM(roll.Campaign, user)

	if !isRoller && !isGM {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	if err := initializers.DB.Delete(&roll).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to destroy roll"})
		return
	}

	c.JSON(200, serializers.SerializeRoll(*roll))
}
