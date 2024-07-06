package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
)

func WorldCreate(c *gin.Context) {
	world := helpers.BodyToWorld(c)

	if result := initializers.DB.Create(&world); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create world"})
		return
	}

	c.JSON(200, serializers.SerializeWorld(*world))
}

func WorldIndex(c *gin.Context) {
	var worlds []models.World
	user := helpers.GetUserFromContext(c, false)

	if user != nil {
		initializers.DB.
			Preload("Creators").
			Where("public = ? OR id in (SELECT world_id FROM world_creators WHERE user_id = ?)", true, user.ID).
			Find(&worlds)
	} else {
		initializers.DB.
			Preload("Creators").
			Where("Public = ?", true).
			Find(&worlds)
	}

	c.JSON(200, gin.H{
		"worlds": serializers.SerializeWorlds(worlds),
	})
}

func WorldRetrieve(c *gin.Context) {
	world := helpers.GetWorldFromSlug(c)
	user := helpers.GetUserFromContext(c, false)
	allowed := helpers.HasWorldAccess(world, user)

	if !allowed {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeWorld(*world))
}

func WorldUpdate(c *gin.Context) {
	world := helpers.GetWorldFromSlug(c)
	user := helpers.GetUserFromContext(c, false)
	isCreator := helpers.IsWorldCreator(world, user)

	if !isCreator {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	newWorld := helpers.BodyToWorld(c)
	world.Name = newWorld.Name
	world.Slug = newWorld.Slug
	world.Public = newWorld.Public
	world.Creators = newWorld.Creators

	if err := initializers.DB.Save(world).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update world"})
		return
	}

	c.JSON(200, serializers.SerializeWorld(*world))
}
