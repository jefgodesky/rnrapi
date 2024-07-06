package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
	"gorm.io/gorm"
	"strings"
)

func getUserByUsername(c *gin.Context) (*models.User, bool) {
	username := c.Param("username")
	var user models.User
	result := initializers.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": fmt.Sprintf("User %s not found", username)})
			return nil, false
		}
		c.JSON(500, gin.H{"error": "Failed to retrieve user"})
		return nil, false
	}

	return &user, true
}

func UserCreate(c *gin.Context) {
	var body struct {
		Username string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	key, err := models.GenerateAPIKey()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to generate API key"})
		return
	}

	hash, err := models.HashAPIKey(key)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to hash API key"})
		return
	}

	user := models.User{Username: body.Username, APIKey: hash, Active: true}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		username := "duplicate key value violates unique constraint"
		if strings.Contains(result.Error.Error(), username) {
			c.JSON(409, gin.H{"error": fmt.Sprintf("Username %s already exists", body.Username)})
		} else {
			c.JSON(400, gin.H{"error": "Failed to create user"})
		}
		return
	}

	location := fmt.Sprintf("/v1/users/%s", user.Username)
	c.Header("Location", location)
	c.JSON(200, serializers.SerializeUser(user, &key))
}

func UserIndex(c *gin.Context) {
	var users []models.User
	initializers.DB.Where("active = ?", true).Find(&users)
	c.JSON(200, gin.H{
		"users": serializers.SerializeUsers(users),
	})
}

func UserRetrieve(c *gin.Context) {
	user, ok := getUserByUsername(c)
	if !ok {
		return
	}
	c.JSON(200, serializers.SerializeUser(*user, nil))
}

func UserUpdate(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	apiKey := c.GetHeader("Authorization")
	if apiKey == "" || !strings.HasPrefix(apiKey, "Bearer ") {
		c.JSON(401, gin.H{"error": "API key required"})
		return
	}
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")

	user, ok := getUserByUsername(c)
	if !ok {
		return
	}

	if err := models.CheckAPIKey(apiKey, user.APIKey); err != nil {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	user.Username = body.Username
	if err := initializers.DB.Save(user).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update user"})
		return
	}

	c.JSON(200, serializers.SerializeUser(*user, nil))
}
