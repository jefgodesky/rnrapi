package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/models"
)

func GetWorldFromSlug(c *gin.Context) *models.World {
	slug := c.Param("slug")
	return GetWorld(c, slug)
}

func GetCampaignFromSlug(c *gin.Context) *models.Campaign {
	world := c.Param("world")
	slug := c.Param("slug")
	return GetCampaign(c, world, slug)
}
