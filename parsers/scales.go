package parsers

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/models"
)

func BodyToScale(c *gin.Context) *models.Scale {
	type LevelBody struct {
		Order       int    `json:"order"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var body struct {
		Name        string      `json:"name"`
		Slug        string      `json:"slug"`
		Description string      `json:"description"`
		Levels      []LevelBody `json:"levels"`
		Public      *bool       `json:"public"`
		Author      *string     `json:"author"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	scaleSlug := body.Slug
	if scaleSlug == "" {
		scaleSlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	author := helpers.GetUserFromContext(c, true)
	if body.Author != nil {
		author = helpers.GetUser(c, *body.Author, false)
	}

	var levels = make([]models.Level, 0)
	for _, level := range body.Levels {
		levels = append(levels, models.Level{
			Order:       level.Order,
			Name:        level.Name,
			Description: level.Description,
		})
	}

	scale := models.Scale{
		Name:        body.Name,
		Slug:        scaleSlug,
		Description: body.Description,
		Levels:      levels,
		Public:      isPublic,
		Author:      *author,
	}

	return &scale
}
