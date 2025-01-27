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

func GetUserFromUsername(c *gin.Context) *models.User {
	username := c.Param("username")
	return GetUser(c, username, true)
}

func GetCharacterFromID(c *gin.Context) *models.Character {
	id := c.Param("id")
	return GetCharacter(c, id)
}

func GetScaleFromSlug(c *gin.Context) *models.Scale {
	slug := c.Param("slug")
	return GetScale(c, slug)
}

func GetScrollFromID(c *gin.Context) *models.Scroll {
	id := c.Param("id")
	return GetScroll(c, id)
}

func GetTableFromSlug(c *gin.Context) *models.Table {
	slug := c.Param("slug")
	return GetTable(c, slug)
}

func GetRollFromID(c *gin.Context) *models.Roll {
	id := c.Param("id")
	return GetRoll(c, id)
}

func GetKeyFromID(c *gin.Context) *models.Key {
	id := c.Param("id")
	return GetKey(c, id)
}

func GetEmailFromID(c *gin.Context) *models.Email {
	id := c.Param("id")
	return GetEmail(c, id)
}
