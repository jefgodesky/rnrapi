package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
)

func WorldCreate(c *gin.Context) {
	var body struct {
		Name     string   `json:"name"`
		Slug     string   `json:"slug"`
		Public   *bool    `json:"public"`
		Creators []string `json:"creators"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	worldSlug := body.Slug
	if worldSlug == "" {
		worldSlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	var creators []models.User
	for _, creator := range body.Creators {
		var user models.User
		result := initializers.DB.Where("username = ? AND active = ?", creator, true).First(&user)
		if result.Error == nil {
			creators = append(creators, user)
		}
	}

	if len(creators) == 0 {
		authUser := helpers.GetUserFromContext(c)
		creators = append(creators, *authUser)
	}

	world := models.World{
		Name:     body.Name,
		Slug:     worldSlug,
		Public:   isPublic,
		Creators: creators,
	}

	if result := initializers.DB.Create(&world); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create world"})
		return
	}

	c.JSON(200, serializers.SerializeWorld(world))
}
