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

func BodyToCampaign(c *gin.Context) *models.Campaign {
	var body struct {
		Slug        string   `json:"slug"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		GMs         []string `json:"creators"`
		Public      *bool    `json:"public"`
		World       string   `json:"world"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return nil
	}

	campaignSlug := body.Slug
	if campaignSlug == "" {
		campaignSlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	world := GetWorld(c, body.World)
	if world == nil {
		c.JSON(400, gin.H{"error": "World not found"})
		return nil
	}

	var gms []models.User
	for _, gm := range body.GMs {
		var user models.User
		result := initializers.DB.
			Where("username = ? AND active = ?", gm, true).
			First(&user)
		if result.Error == nil {
			gms = append(gms, user)
		}
	}

	if len(gms) == 0 {
		authUser := GetUserFromContext(c, true)
		gms = append(gms, *authUser)
	}

	campaign := models.Campaign{
		Slug:        campaignSlug,
		Name:        body.Name,
		Description: body.Description,
		GMs:         gms,
		Public:      isPublic,
		WorldID:     world.ID,
		World:       *world,
	}

	return &campaign
}
