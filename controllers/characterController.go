package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
)

func CharacterCreate(c *gin.Context) {
	character := helpers.BodyToCharacter(c)
	if character == nil {
		return
	}

	if result := initializers.DB.Create(&character); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create character"})
		return
	}

	c.JSON(200, serializers.SerializeCharacter(*character))
}

func CharacterIndex(c *gin.Context) {
	var characters []models.Character
	user := helpers.GetUserFromContext(c, false)

	if user != nil {
		initializers.DB.
			Where("public = ? OR player_id = ?", true, user.ID).
			Find(&characters)
	} else {
		initializers.DB.
			Where("public = ?", true).
			Find(&characters)
	}

	c.JSON(200, gin.H{
		"characters": serializers.SerializeCharacters(characters),
	})
}

func CharacterRetrieve(c *gin.Context) {
	char := helpers.GetCharacterFromSlug(c)
	user := helpers.GetUserFromContext(c, false)
	if char == nil {
		return
	}

	if !char.Public && char.PlayerID != user.ID {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeCharacter(*char))
}

func CharacterUpdate(c *gin.Context) {
	char := helpers.CharacterPlayerOnly(c)
	if char == nil {
		return
	}

	newChar := helpers.BodyToCharacter(c)
	char.Name = newChar.Name
	char.Description = newChar.Description
	char.Str = newChar.Str
	char.Dex = newChar.Dex
	char.Con = newChar.Con
	char.Int = newChar.Int
	char.Wis = newChar.Wis
	char.Cha = newChar.Cha
	char.Notes = newChar.Notes
	char.Public = newChar.Public
	char.PlayerID = newChar.PlayerID
	char.Player = newChar.Player

	if err := initializers.DB.Save(char).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update character"})
		return
	}

	c.JSON(200, serializers.SerializeCharacter(*char))
}
