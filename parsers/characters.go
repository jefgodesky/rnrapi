package parsers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
)

func IdsToCharacters(ids []string) []models.Character {
	var characters []models.Character
	for _, id := range ids {
		var char models.Character
		result := initializers.DB.First(&char, "id = ?", id)
		if result.Error == nil {
			characters = append(characters, char)
		}
	}
	return characters
}

func BodyToCharacter(c *gin.Context) *models.Character {
	type AbilitiesBody struct {
		Strength     int `json:"strength"`
		Dexterity    int `json:"dexterity"`
		Constitution int `json:"constitution"`
		Intelligence int `json:"intelligence"`
		Wisdom       int `json:"wisdom"`
		Charisma     int `json:"charisma"`
	}

	var body struct {
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Abilities   AbilitiesBody `json:"abilities"`
		Notes       string        `json:"notes"`
		Public      *bool         `json:"public"`
		PC          bool          `json:"pc"`
		Player      string        `json:"player"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	player := helpers.GetUser(c, body.Player, false)
	if player == nil {
		player = helpers.GetUserFromContext(c, false)
	}

	char := models.Character{
		Name:        body.Name,
		Description: body.Description,
		Str:         body.Abilities.Strength,
		Dex:         body.Abilities.Dexterity,
		Con:         body.Abilities.Constitution,
		Int:         body.Abilities.Intelligence,
		Wis:         body.Abilities.Wisdom,
		Cha:         body.Abilities.Charisma,
		Notes:       body.Notes,
		PC:          body.PC,
		Public:      isPublic,
		PlayerID:    player.ID,
		Player:      *player,
	}

	return &char
}
