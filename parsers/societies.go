package parsers

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jefgodesky/rnrapi/enums"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/models"
)

func BodyToSociety(c *gin.Context) *models.Society {
	var body struct {
		Slug        string           `json:"slug"`
		Name        string           `json:"name"`
		Description string           `json:"description"`
		Favored     [2]enums.Ability `json:"favored"`
		Languages   string           `json:"languages"`
		Public      *bool            `json:"public"`
		World       string           `json:"world"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	societySlug := body.Slug
	if societySlug == "" {
		societySlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	world := helpers.GetWorld(c, body.World)
	if world == nil {
		c.JSON(400, gin.H{"error": "World not found"})
		c.Abort()
		return nil
	}

	society := models.Society{
		Slug:        societySlug,
		Name:        body.Name,
		Description: body.Description,
		Favored:     enums.AbilityPair{body.Favored[0], body.Favored[1]},
		Languages:   body.Languages,
		Public:      isPublic,
		WorldID:     world.ID,
		World:       *world,
	}

	return &society
}
