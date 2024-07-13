package parsers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/models"
	"strings"
)

func BodyToRoll(c *gin.Context) *models.Roll {
	var body struct {
		Note      *string `json:"note"`
		Table     string  `json:"table"`
		Roller    *string `json:"roller"`
		Character *string `json:"character"`
		Ability   *string `json:"ability"`
		Campaign  *string `json:"campaign"`
		Modifier  *int    `json:"modifier"`
	}

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	var roller *models.User = nil
	if body.Roller != nil {
		roller = helpers.GetUser(c, *body.Roller, false)
	}

	if roller == nil {
		roller = helpers.GetUserFromContext(c, false)
	}

	table := helpers.GetTable(c, body.Table)
	if table == nil {
		c.Abort()
		return nil
	}

	modifier := 0
	if body.Modifier != nil {
		modifier = *body.Modifier
	}

	roll := models.Roll{
		Table:    *table,
		Modifier: modifier,
		Log:      "",
		Results:  "",
	}

	if body.Note != nil {
		roll.Note = body.Note
	}

	if roller != nil {
		roll.Roller = roller
	}

	if body.Character != nil {
		char := helpers.GetCharacter(c, *body.Character)
		if char != nil {
			roll.Character = char
		}
	}

	if body.Ability != nil {
		roll.Ability = body.Ability
	}

	if body.Campaign != nil {
		parts := strings.Split(*body.Campaign, "/")
		if len(parts) > 1 {
			campaign := helpers.GetCampaign(c, parts[0], parts[1])
			if campaign != nil {
				roll.Campaign = campaign
			}
		}
	}

	return &roll
}
