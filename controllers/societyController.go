package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
	"gorm.io/gorm/clause"
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

func SocietyIndex(c *gin.Context) {
	var societies []models.Society
	initializers.DB.Preload(clause.Associations).Find(&societies)

	user := helpers.GetUserFromContext(c, false)
	filtered := helpers.FilterSocietiesWorldAccess(societies, user)
	c.JSON(200, gin.H{
		"societies": serializers.SerializeSocieties(filtered),
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

	newSociety := helpers.BodyToSociety(c)
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
