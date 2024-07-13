package parsers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/models"
)

func BodyToEmail(c *gin.Context) *models.Email {
	var body struct {
		Address string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	authUser := helpers.GetUserFromContext(c, true)

	email := models.Email{
		Address: body.Address,
		User:    *authUser,
	}

	models.SetVerificationCode(&email)

	return &email
}
