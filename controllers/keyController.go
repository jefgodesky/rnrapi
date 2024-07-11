package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
)

func KeyCreate(c *gin.Context) {
	username, password, label, ephemeral := helpers.BodyToKeyRequest(c)
	if username == "" || password == "" {
		return
	}

	user := helpers.GetUser(c, username, true)
	if !user.Active {
		c.JSON(403, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}

	token, hash, plaintext := helpers.GenerateAPIKey(c)
	if token == "" || hash == "" || plaintext == "" {
		c.JSON(400, gin.H{"error": "Failed to create API key"})
		c.Abort()
		return
	}

	key := models.Key{Token: token, Secret: hash, User: *user, Label: label, Ephemeral: ephemeral}
	result := initializers.DB.Create(&key)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Failed to create API key"})
		c.Abort()
		return
	}

	c.JSON(200, serializers.SerializeKey(key, &plaintext))
}

func KeyIndex(c *gin.Context) {
	user := helpers.GetUserFromContext(c, true)
	var keys []models.Key
	query := initializers.DB.Model(&models.Key{}).Where("user_id = ?", user.ID)

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&keys)
	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"keys":      serializers.SerializeKeys(keys),
	})
}

func KeyRetrieve(c *gin.Context) {
	key := helpers.GetKeyFromID(c)
	c.JSON(200, serializers.SerializeKey(*key, nil))
}

func KeyDestroy(c *gin.Context) {
	key := helpers.GetKeyFromID(c)
	if err := initializers.DB.Delete(key).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to delete key"})
		c.Abort()
		return
	}

	c.Status(204)
}
