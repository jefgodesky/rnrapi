package helpers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

var PreloadPaths = map[string][]string{
	"World":    {"Creators"},
	"Campaign": {"GMs", "PCs", "World", "World.Creators"},
	"Species":  {"World", "World.Creators"},
	"Society":  {"World", "World.Creators"},
	"Table":    {"Rows", "Author"},
}

func GetInstance(c *gin.Context, model interface{}, slug string, conditions map[string]interface{}) bool {
	modelType := reflect.TypeOf(model).Elem().Name()
	preloadPaths, ok := PreloadPaths[modelType]
	if !ok {
		c.JSON(500, gin.H{"error": fmt.Sprintf("No preload path for model type %s", modelType)})
		return false
	}

	db := initializers.DB
	for _, path := range preloadPaths {
		db = db.Preload(path)
	}

	query := db.Where("slug = ?", slug)
	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	result := query.First(model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": fmt.Sprintf("%s %s not found", modelType, slug)})
			return false
		}
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to retrieve %s", modelType)})
		return false
	}

	return true
}

func GetWorld(c *gin.Context, slug string) *models.World {
	var world models.World
	if !GetInstance(c, &world, slug, map[string]interface{}{}) {
		return nil
	}
	return &world
}

func GetCampaign(c *gin.Context, world string, slug string) *models.Campaign {
	w := GetWorld(c, world)
	if w == nil {
		return nil
	}

	var campaign models.Campaign
	if !GetInstance(c, &campaign, slug, map[string]interface{}{
		"world_id": w.ID,
	}) {
		return nil
	}
	return &campaign
}

func GetSpecies(c *gin.Context, world string, slug string) *models.Species {
	w := GetWorld(c, world)
	if w == nil {
		return nil
	}

	var species models.Species
	if !GetInstance(c, &species, slug, map[string]interface{}{
		"world_id": w.ID,
	}) {
		return nil
	}
	return &species
}

func GetSociety(c *gin.Context, world string, slug string) *models.Society {
	w := GetWorld(c, world)
	if w == nil {
		return nil
	}

	var society models.Society
	if !GetInstance(c, &society, slug, map[string]interface{}{
		"world_id": w.ID,
	}) {
		return nil
	}
	return &society
}

func GetUser(c *gin.Context, username string, required bool) *models.User {
	var user models.User
	result := initializers.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if required {
				c.JSON(404, gin.H{"error": fmt.Sprintf("User %s not found", username)})
				c.Abort()
			}
			return nil
		}
		if required {
			c.JSON(500, gin.H{"error": "Failed to retrieve user"})
		}
		return nil
	}

	return &user
}

func GetCharacter(c *gin.Context, id string) *models.Character {
	var char models.Character
	result := initializers.DB.Preload(clause.Associations).Where("id = ?", id).First(&char)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": fmt.Sprintf("Character %s not found", id)})
			return nil
		}
		c.JSON(500, gin.H{"error": "Failed to retrieve character"})
		return nil
	}

	return &char
}

func GetScroll(c *gin.Context, id string) *models.Scroll {
	var scroll models.Scroll
	result := initializers.DB.
		Preload(clause.Associations).
		Preload("Campaign.World").
		Where("id = ?", id).First(&scroll)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": fmt.Sprintf("Scroll %s not found", id)})
			return nil
		}
		c.JSON(500, gin.H{"error": "Failed to retrieve scroll"})
		return nil
	}

	return &scroll
}

func GetTable(c *gin.Context, slug string) *models.Table {
	var table models.Table
	if !GetInstance(c, &table, slug, map[string]interface{}{}) {
		return nil
	}
	return &table
}

func GetRoll(c *gin.Context, id string) *models.Roll {
	var roll models.Roll
	result := initializers.DB.
		Preload(clause.Associations).
		Preload("Campaign.World").
		Preload("Table.Author").
		Where("id = ?", id).First(&roll)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": fmt.Sprintf("Roll %s not found", id)})
			return nil
		}
		c.JSON(500, gin.H{"error": "Failed to retrieve roll"})
		return nil
	}

	return &roll
}
