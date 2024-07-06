package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
)

func BodyToUserFields(c *gin.Context) (string, string, string) {
	var body struct {
		Username string
		Name     string
		Bio      string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return "", "", ""
	}

	return body.Username, body.Name, body.Bio
}

func BodyToWorld(c *gin.Context) *models.World {
	var body struct {
		Name     string   `json:"name"`
		Slug     string   `json:"slug"`
		Public   *bool    `json:"public"`
		Creators []string `json:"creators"`
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

	var creators []models.User
	for _, creator := range body.Creators {
		var user models.User
		result := initializers.DB.
			Where("username = ? AND active = ?", creator, true).
			First(&user)
		if result.Error == nil {
			creators = append(creators, user)
		}
	}

	if len(creators) == 0 {
		authUser := GetUserFromContext(c, true)
		creators = append(creators, *authUser)
	}

	world := models.World{
		Name:     body.Name,
		Slug:     worldSlug,
		Public:   isPublic,
		Creators: creators,
	}

	return &world
}
