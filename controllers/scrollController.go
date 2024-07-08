package controllers

import (
	"fmt"
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
	scroll := helpers.GetScrollFromID(c)
	user := helpers.GetUserFromContext(c, false)
	allowed := scroll.Public || helpers.IsScrollReader(scroll, user)

	if !allowed {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeScroll(*scroll))
}

func ScrollUpdate(c *gin.Context) {
	scroll := helpers.ScrollWriterOnly(c)
	if scroll == nil {
		return
	}

	newScroll := helpers.BodyToScroll(c)
	scroll.Name = newScroll.Name
	scroll.Description = newScroll.Description
	scroll.Seals = newScroll.Seals
	scroll.Writers = newScroll.Writers
	scroll.Readers = newScroll.Readers
	scroll.Public = newScroll.Public
	scroll.CampaignID = newScroll.CampaignID
	scroll.Campaign = newScroll.Campaign

	if err := initializers.DB.Save(scroll).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update scroll"})
		return
	}

	c.JSON(200, serializers.SerializeScroll(*scroll))
}

func ScrollDestroy(c *gin.Context) {
	scroll := helpers.ScrollWriterOnly(c)
	if scroll == nil {
		return
	}

	if err := initializers.DB.Delete(&scroll).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to destroy scroll"})
		return
	}

	c.JSON(200, serializers.SerializeScroll(*scroll))
}

func ScrollSeal(c *gin.Context) {
	scroll := helpers.ScrollWriterOnly(c)
	if scroll == nil {
		return
	}

	scroll.Seals += 1

	if err := initializers.DB.Save(scroll).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update scroll"})
		return
	}

	c.JSON(200, serializers.SerializeScroll(*scroll))
}

func ScrollUnseal(c *gin.Context) {
	scroll := helpers.ScrollWriterOnly(c)
	if scroll == nil {
		return
	}

	if scroll.Seals > 0 {
		scroll.Seals -= 1
	}

	if err := initializers.DB.Save(scroll).Error; err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"Error": "Failed to update scroll"})
		return
	}

	c.JSON(200, serializers.SerializeScroll(*scroll))
}
