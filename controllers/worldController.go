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

func WorldCreate(c *gin.Context) {
	world := parsers.BodyToWorld(c)
	if world == nil {
		return
	}

	if result := initializers.DB.Create(&world); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create world"})
		return
	}

	c.JSON(200, serializers.SerializeWorld(*world))
}

func WorldIndex(c *gin.Context) {
	var worlds []models.World
	user := helpers.GetUserFromContext(c, false)
	query := initializers.DB.Model(&models.World{}).Preload(clause.Associations)

	if user != nil {
		query.Where("public = ? OR id in (SELECT world_id FROM world_creators WHERE user_id = ?)", true, user.ID)
	} else {
		query.Where("public = ?", true)
	}

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&worlds)

	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"worlds":    serializers.SerializeWorlds(worlds),
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
	world := helpers.WorldCreatorOnly(c)
	if world == nil {
		return
	}

	newWorld := parsers.BodyToWorld(c)
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

func WorldDestroy(c *gin.Context) {
	world := helpers.WorldCreatorOnly(c)
	if world == nil {
		return
	}

	if err := initializers.DB.Delete(&world).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to destroy world"})
		return
	}

	c.JSON(200, serializers.SerializeWorld(*world))
}
