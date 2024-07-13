package parsers

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/models"
)

func BodyToCampaign(c *gin.Context) *models.Campaign {
	var body struct {
		Slug        string   `json:"slug"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		GMs         []string `json:"gms"`
		PCs         []string `json:"pcs"`
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

	world := helpers.GetWorld(c, body.World)
	if world == nil {
		c.JSON(400, gin.H{"error": "World not found"})
		return nil
	}

	authUser := helpers.GetUserFromContext(c, true)
	gms := UsernamesToUsersWithDefault(body.GMs, *authUser)
	pcs := IdsToCharacters(body.PCs)

	campaign := models.Campaign{
		Slug:        campaignSlug,
		Name:        body.Name,
		Description: body.Description,
		GMs:         gms,
		PCs:         pcs,
		Public:      isPublic,
		WorldID:     world.ID,
		World:       *world,
	}

	return &campaign
}
