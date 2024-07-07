package helpers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"gorm.io/gorm"
	"reflect"
)

var PreloadPaths = map[string][]string{
	"World":    {"Creators"},
	"Campaign": {"GMs", "World", "World.Creators"},
	"Species":  {"World", "World.Creators"},
	"Society":  {"World", "World.Creators"},
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
			c.JSON(404, gin.H{"error": fmt.Sprintf("%T %s not found", model, slug)})
			return false
		}
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to retrieve %T", model)})
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

func GetUser(c *gin.Context, username string) *models.User {
	var user models.User
	result := initializers.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": fmt.Sprintf("User %s not found", username)})
			return nil
		}
		c.JSON(500, gin.H{"error": "Failed to retrieve user"})
		return nil
	}

	return &user
}
