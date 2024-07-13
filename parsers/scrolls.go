package parsers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/models"
	"strings"
)

func BodyToScroll(c *gin.Context) *models.Scroll {
	var body struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Seals       uint     `json:"seals"`
		Writers     []string `json:"writers"`
		Readers     []string `json:"readers"`
		Public      *bool    `json:"public"`
		Campaign    *string  `json:"campaign"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return nil
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	authUser := helpers.GetUserFromContext(c, true)
	writers := UsernamesToUsersWithDefault(body.Writers, *authUser)
	readers := UsernamesToUsersWithDefault(body.Readers, *authUser)

	var campaign *models.Campaign
	if body.Campaign != nil {
		parts := strings.Split(*body.Campaign, "/")
		if len(parts) > 1 {
			campaign = helpers.GetCampaign(c, parts[0], parts[1])
		}
	}

	scroll := models.Scroll{
		Name:        body.Name,
		Description: body.Description,
		Seals:       body.Seals,
		Writers:     writers,
		Readers:     readers,
		Public:      isPublic,
		Campaign:    campaign,
	}

	return &scroll
}
