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

func GetSpeciesFromSlug(c *gin.Context) *models.Species {
	world := c.Param("world")
	slug := c.Param("slug")
	return GetSpecies(c, world, slug)
}

func GetSocietyFromSlug(c *gin.Context) *models.Society {
	world := c.Param("world")
	slug := c.Param("slug")
	return GetSociety(c, world, slug)
}
