package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/serializers"
)

func SpeciesCreate(c *gin.Context) {
	species := helpers.BodyToSpecies(c)

	if result := initializers.DB.Create(&species); result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(500, gin.H{"error": "Failed to create species"})
		return
	}

	c.JSON(200, serializers.SerializeSpecies(*species))
}
