package helpers

import (
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
