package parsers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jefgodesky/rnrapi/enums"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/models"
)

func BodyToSpecies(c *gin.Context) *models.Species {
	var body struct {
		Slug        string           `json:"slug"`
		Name        string           `json:"name"`
		Description string           `json:"description"`
		Affinities  [2]enums.Ability `json:"affinities"`
		Aversion    enums.Ability    `json:"aversion"`
		Stages      json.RawMessage  `json:"stages"`
		Public      *bool            `json:"public"`
		World       string           `json:"world"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	speciesSlug := body.Slug
	if speciesSlug == "" {
		speciesSlug = slug.Make(body.Name)
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

	species := models.Species{
		Slug:        speciesSlug,
		Name:        body.Name,
		Description: body.Description,
		Affinities:  enums.AbilityPair{body.Affinities[0], body.Affinities[1]},
		Aversion:    body.Aversion,
		Stages:      body.Stages,
		Public:      isPublic,
		WorldID:     world.ID,
		World:       *world,
	}

	return &species
}
