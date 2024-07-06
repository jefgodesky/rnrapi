package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/models"
)

func GetWorldFromSlug(c *gin.Context) *models.World {
	slug := c.Param("slug")
	return GetWorld(slug, c)
}
