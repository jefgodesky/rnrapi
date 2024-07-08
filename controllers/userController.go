package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/serializers"
	"strings"
)

func UserCreate(c *gin.Context) {
	username, name, bio := helpers.BodyToUserFields(c)
	if username == "" || name == "" || bio == "" {
		return
	}

	token, hash, key := helpers.GenerateAPIKey(c)
	if token == "" || hash == "" || key == "" {
		return
	}

	user := models.User{Username: username, Name: name, Bio: bio, Token: token, Secret: hash, Active: true}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		usernameError := "duplicate key value violates unique constraint"
		if strings.Contains(result.Error.Error(), usernameError) {
			c.JSON(409, gin.H{"error": fmt.Sprintf("Username %s already exists", username)})
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
	query := initializers.DB.Model(&models.User{}).Where("active = ?", true)
	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&users)
	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"users":     serializers.SerializeUsers(users),
	})
}

func UserRetrieve(c *gin.Context) {
	user := helpers.GetUserFromUsername(c)
	if user == nil {
		return
	}
	c.JSON(200, serializers.SerializeUser(*user, nil))
}

func UserUpdate(c *gin.Context) {
	username, name, bio := helpers.BodyToUserFields(c)
	if username == "" || name == "" || bio == "" {
		return
	}

	user := helpers.GetUserFromContext(c, true)
	user.Username = username
	user.Name = name
	user.Bio = bio

	if err := initializers.DB.Save(user).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update user"})
		return
	}

	c.JSON(200, serializers.SerializeUser(*user, nil))
}
