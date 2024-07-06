package helpers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"gorm.io/gorm"
)

func GetWorld(slug string, c *gin.Context) *models.World {
	var world models.World
	result := initializers.DB.
		Preload("Creators").
		Where("slug = ?", slug).
		First(&world)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": fmt.Sprintf("World %s not found", slug)})
			return nil
		}
		c.JSON(500, gin.H{"error": "Failed to retrieve world"})
		return nil
	}

	return &world
}
