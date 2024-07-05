package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
)

func UserCreate(c *gin.Context) {
	var body struct {
		Username string
		Password string
	}

	c.Bind(&body)

	user := models.User{Username: body.Username, Password: body.Password, Active: true}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, user)
}
