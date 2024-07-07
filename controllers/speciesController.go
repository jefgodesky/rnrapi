package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
	"gorm.io/gorm/clause"
)

func SpeciesCreate(c *gin.Context) {
	species := helpers.BodyToSpecies(c)

	if result := initializers.DB.Create(&species); result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(500, gin.H{"error": "Failed to create species"})
		return
	}

	c.JSON(200, serializers.SerializeSpecies(*species))
}

func SpeciesIndex(c *gin.Context) {
	var species []models.Species
	initializers.DB.Preload(clause.Associations).Find(&species)

	user := helpers.GetUserFromContext(c, false)
	filtered := helpers.FilterSpeciesWorldAccess(species, user)
	c.JSON(200, gin.H{
		"species": serializers.SerializeSpp(filtered),
	})
}

func SpeciesRetrieve(c *gin.Context) {
	species := helpers.GetSpeciesFromSlug(c)
	user := helpers.GetUserFromContext(c, false)
	allowed := helpers.HasWorldAccess(&species.World, user)

	if !allowed {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeSpecies(*species))
}

func SpeciesUpdate(c *gin.Context) {
	species := helpers.SpeciesCreatorOnly(c)
	if species == nil {
		return
	}

	newSpecies := helpers.BodyToSpecies(c)
	species.Slug = newSpecies.Slug
	species.Name = newSpecies.Name
	species.Description = newSpecies.Description
	species.Affinities = newSpecies.Affinities
	species.Aversion = newSpecies.Aversion
	species.Public = newSpecies.Public
	species.WorldID = newSpecies.WorldID
	species.World = newSpecies.World

	if err := initializers.DB.Save(species).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update species"})
		return
	}

	c.JSON(200, serializers.SerializeSpecies(*species))
}

func SpeciesDestroy(c *gin.Context) {
	species := helpers.SpeciesCreatorOnly(c)
	if species == nil {
		return
	}

	if err := initializers.DB.Delete(&species).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to destroy species"})
		return
	}

	c.JSON(200, serializers.SerializeSpecies(*species))
}
