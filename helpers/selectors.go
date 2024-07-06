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
}

func GetInstance(c *gin.Context, model interface{}, slug string) bool {
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

	result := db.Where("slug = ?", slug).First(model)
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
	if !GetInstance(c, &world, slug) {
		return nil
	}
	return &world
}

func GetCampaign(c *gin.Context, slug string) *models.Campaign {
	var campaign models.Campaign
	if !GetInstance(c, &campaign, slug) {
		return nil
	}
	return &campaign
}
