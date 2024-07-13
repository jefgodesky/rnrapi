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

func SocietyCreate(c *gin.Context) {
	society := parsers.BodyToSociety(c)
	if society == nil {
		return
	}

	if result := initializers.DB.Create(&society); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create society"})
		return
	}

	c.JSON(200, serializers.SerializeSociety(*society))
}

func SocietyIndex(c *gin.Context) {
	var societies []models.Society
	user := helpers.GetUserFromContext(c, false)
	query := initializers.DB.
		Model(&models.Society{}).
		Preload(clause.Associations)

	if user != nil {
		query.Joins("JOIN worlds ON worlds.id = societies.world_id").
			Where("(societies.public = ? AND worlds.public = ?) OR societies.world_id IN (SELECT world_id FROM world_creators WHERE user_id = ?)", true, true, user.ID)
	} else {
		query.Joins("JOIN worlds ON worlds.id = societies.world_id").
			Where("societies.public = ? AND worlds.public = ?", true, true)
	}

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&societies)

	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"societies": serializers.SerializeSocieties(societies),
	})
}

func SocietyRetrieve(c *gin.Context) {
	society := helpers.GetSocietyFromSlug(c)
	user := helpers.GetUserFromContext(c, false)
	allowed := helpers.HasWorldAccess(&society.World, user)

	if !allowed {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeSociety(*society))
}

func SocietyUpdate(c *gin.Context) {
	society := helpers.SocietyCreatorOnly(c)
	if society == nil {
		return
	}

	newSociety := parsers.BodyToSociety(c)
	society.Slug = newSociety.Slug
	society.Name = newSociety.Name
	society.Description = newSociety.Description
	society.Favored = newSociety.Favored
	society.Languages = newSociety.Languages
	society.Public = newSociety.Public
	society.WorldID = newSociety.WorldID
	society.World = newSociety.World

	if err := initializers.DB.Save(society).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update society"})
		return
	}

	c.JSON(200, serializers.SerializeSociety(*society))
}

func SocietyDestroy(c *gin.Context) {
	society := helpers.SocietyCreatorOnly(c)
	if society == nil {
		return
	}

	if err := initializers.DB.Delete(&society).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to destroy society"})
		return
	}

	c.JSON(200, serializers.SerializeSociety(*society))
}
