package parsers

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/models"
)

func BodyToWorld(c *gin.Context) *models.World {
	var body struct {
		Name        string   `json:"name"`
		Slug        string   `json:"slug"`
		Description string   `json:"description"`
		Public      *bool    `json:"public"`
		Creators    []string `json:"creators"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return nil
	}

	worldSlug := body.Slug
	if worldSlug == "" {
		worldSlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	authUser := helpers.GetUserFromContext(c, true)
	creators := UsernamesToUsersWithDefault(body.Creators, *authUser)

	world := models.World{
		Name:        body.Name,
		Slug:        worldSlug,
		Description: body.Description,
		Public:      isPublic,
		Creators:    creators,
	}

	return &world
}
