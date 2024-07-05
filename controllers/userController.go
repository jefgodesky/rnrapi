package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"strings"
)

func UserCreate(c *gin.Context) {
	var body struct {
		Username string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	hash, err := models.HashPassword(body.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{Username: body.Username, Password: hash, Active: true}
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
	c.Status(201)
}
