package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/serializers"
)

func RollCreate(c *gin.Context) {
	roll := helpers.BodyToRoll(c)
	if roll == nil {
		return
	}

	helpers.RollOnTable(roll.Table, roll, roll.Modifier)

	if result := initializers.DB.Create(&roll); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create roll record"})
		return
	}

	c.JSON(200, serializers.SerializeRoll(*roll))
}
