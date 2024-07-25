package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/parsers"
	"github.com/jefgodesky/rnrapi/serializers"
	"gorm.io/gorm/clause"
)

func SpeciesCreate(c *gin.Context) {
	species := parsers.BodyToSpecies(c)
	if species == nil {
		return
	}

	if result := initializers.DB.Create(species); result.Error != nil {
		fmt.Println(species)
		fmt.Println(result.Error)
		c.JSON(500, gin.H{"error": "Failed to create species"})
		return
	}

	c.JSON(200, serializers.SerializeSpecies(*species))
}

func SpeciesIndex(c *gin.Context) {
	var species []models.Species
	user := helpers.GetUserFromContext(c, false)
	query := initializers.DB.
		Model(&models.Species{}).
		Preload(clause.Associations)

	if user != nil {
		query.Joins("JOIN worlds ON worlds.id = species.world_id").
			Where("(species.public = ? AND worlds.public = ?) OR species.world_id IN (SELECT world_id FROM world_creators WHERE user_id = ?)", true, true, user.ID)
	} else {
		query.Joins("JOIN worlds ON worlds.id = species.world_id").
			Where("species.public = ? AND worlds.public = ?", true, true)

	}

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&species)

	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"species":   serializers.SerializeSpp(species),
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

	newSpecies := parsers.BodyToSpecies(c)
	species.Slug = newSpecies.Slug
	species.Name = newSpecies.Name
	species.Description = newSpecies.Description
	species.Affinities = newSpecies.Affinities
	species.Aversion = newSpecies.Aversion
	species.Stages = newSpecies.Stages
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
