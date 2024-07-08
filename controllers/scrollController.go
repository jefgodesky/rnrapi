package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
)

func ScrollCreate(c *gin.Context) {
	scroll := helpers.BodyToScroll(c)

	if result := initializers.DB.Create(&scroll); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create scroll"})
		return
	}

	c.JSON(200, serializers.SerializeScroll(*scroll))
}

func ScrollIndex(c *gin.Context) {
	var scrolls []models.Scroll
	user := helpers.GetUserFromContext(c, false)
	query := initializers.DB.Model(&models.Scroll{})

	if user != nil {
		query.Joins("LEFT JOIN scroll_readers ON scroll_readers.scroll_id = scrolls.id").
			Where("scrolls.public = ? OR scroll_readers.user_id = ?", true, user.ID)
	} else {
		query.Where("public = ?", true)
	}

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&scrolls)

	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"scrolls":   serializers.SerializeScrolls(scrolls),
	})
}

func ScrollRetrieve(c *gin.Context) {
	scroll := helpers.GetScrollFromSlug(c)
	user := helpers.GetUserFromContext(c, false)
	allowed := scroll.Public || helpers.IsScrollReader(scroll, user)

	if !allowed {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeScroll(*scroll))
}