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

func ScaleCreate(c *gin.Context) {
	scale := parsers.BodyToScale(c)

	if result := initializers.DB.Create(&scale); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create scale"})
		return
	}

	c.JSON(200, serializers.SerializeScale(*scale))
}

func ScaleIndex(c *gin.Context) {
	var scales []models.Scale
	user := helpers.GetUserFromContext(c, false)
	query := initializers.DB.
		Preload(clause.Associations).
		Model(&models.Scale{})

	if user != nil {
		query.Where("public = ? OR author_id = ?", true, user.ID)
	} else {
		query.Where("public = ?", true)
	}

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&scales)

	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"scales":    serializers.SerializeScales(scales),
	})
}

func ScaleRetrieve(c *gin.Context) {
	scale := helpers.GetScaleFromSlug(c)
	user := helpers.GetUserFromContext(c, false)
	allowed := scale.Public || scale.Author.ID == user.ID

	if !allowed {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeScale(*scale))
}

func ScaleUpdate(c *gin.Context) {
	scale := helpers.ScaleAuthorOnly(c)
	if scale == nil {
		return
	}

	newScale := parsers.BodyToScale(c)
	scale.Name = newScale.Name
	scale.Slug = newScale.Slug
	scale.Description = newScale.Description
	scale.Levels = newScale.Levels
	scale.Public = newScale.Public
	scale.AuthorID = newScale.AuthorID
	scale.Author = newScale.Author

	if err := initializers.DB.Save(scale).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update scale"})
		return
	}

	c.JSON(200, serializers.SerializeScale(*scale))
}

func ScaleDestroy(c *gin.Context) {
	scale := helpers.ScaleAuthorOnly(c)
	if scale == nil {
		return
	}

	if err := initializers.DB.Delete(&scale).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to destroy scale"})
		return
	}

	c.JSON(200, serializers.SerializeScale(*scale))
}
