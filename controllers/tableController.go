package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/serializers"
)

func TableCreate(c *gin.Context) {
	table := helpers.BodyToTable(c)

	if result := initializers.DB.Create(&table); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create table"})
		return
	}

	c.JSON(200, serializers.SerializeTable(*table))
}
