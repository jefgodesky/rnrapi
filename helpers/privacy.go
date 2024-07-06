package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/models"
)

func IsWorldCreator(world *models.World, user *models.User) bool {
	if user == nil || world == nil {
		return false
	}

	for _, creator := range world.Creators {
		if creator.ID == user.ID {
			return true
		}
	}

	return false
}

func HasWorldAccess(world *models.World, user *models.User) bool {
	if world.Public {
		return true
	}

	if IsWorldCreator(world, user) {
		return true
	}

	return false
}

func WorldCreatorOnly(c *gin.Context) *models.World {
	world := GetWorldFromSlug(c)
	if world == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	isCreator := IsWorldCreator(world, user)
	if !isCreator {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return world
}
