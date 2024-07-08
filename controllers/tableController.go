package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
	"gorm.io/gorm/clause"
)

func TableCreate(c *gin.Context) {
	table := helpers.BodyToTable(c)

	if result := initializers.DB.Create(&table); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create table"})
		return
	}

	c.JSON(200, serializers.SerializeTable(*table))
}

func TableIndex(c *gin.Context) {
	var tables []models.Table
	user := helpers.GetUserFromContext(c, false)
	query := initializers.DB.
		Preload(clause.Associations).
		Model(&models.Table{})

	if user != nil {
		query.Where("public = ? OR author_id = ?", true, user.ID)
	} else {
		query.Where("public = ?", true)
	}

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&tables)

	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"tables":    serializers.SerializeTables(tables),
	})
}

func TableRetrieve(c *gin.Context) {
	table := helpers.GetTableFromSlug(c)
	user := helpers.GetUserFromContext(c, false)
	allowed := table.Public || table.Author.ID == user.ID

	if !allowed {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(200, serializers.SerializeTable(*table))
}
