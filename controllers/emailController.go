package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/parsers"
	"github.com/jefgodesky/rnrapi/serializers"
	"strings"
)

func SendVerificationEmail(email models.Email) error {
	from := "The Ruins & Revolutions Catalogue <catalogue@ruinsandrevolutions.com>"
	to := email.Address
	subject := "Please verify this address"
	fmt.Println(email)

	initializers.DB.Where("id = ?", email.UserID).First(&email.User)
	lines := make([]string, 15)
	lines[0] = fmt.Sprintf("Hello %s,", email.User.Name)
	lines[1] = ""
	lines[2] = "The only thing we will ever use your email address for is to"
	lines[3] = "help you reset your password. If you lose your password or you"
	lines[4] = "can’t remember it, we can create a new one for you and send it"
	lines[5] = "to this email address. This works out great, so long as we can"
	lines[6] = "be sure that when we send that new password to this email"
	lines[7] = "address that we’re sending it to *you* and not some sinister"
	lines[8] = "villain, like a shapeshifting doppleganer, or any other, less"
	lines[9] = "likely threat that may be lurking out there. That’s why we’d"
	lines[10] = "like you to click the link below. That will verify that this"
	lines[11] = "email reached you, so we know where to send your new password"
	lines[12] = "if you should ever need it."
	lines[13] = ""
	lines[14] = fmt.Sprintf("https://ruinsandrevolutions.com/verify/%s,", email.Code)

	return helpers.SendEmail(from, to, subject, strings.Join(lines, "\n"))
}

func EmailCreate(c *gin.Context) {
	email := parsers.BodyToEmail(c)
	if email == nil {
		return
	}

	result := initializers.DB.Create(&email)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Failed to create email"})
		c.Abort()
		return
	}

	err := SendVerificationEmail(*email)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Failed to send verification code email."})
		c.Abort()
		return
	}

	c.JSON(200, serializers.SerializeEmail(*email))
}

func EmailIndex(c *gin.Context) {
	user := helpers.GetUserFromContext(c, true)
	var emails []models.Email
	query := initializers.DB.Model(&models.Email{}).Where("user_id = ?", user.ID)

	var total int64
	query.Count(&total)
	query.Scopes(helpers.Paginate(c)).Find(&emails)
	c.JSON(200, gin.H{
		"total":     total,
		"page":      c.GetInt("page"),
		"page_size": c.GetInt("page_size"),
		"emails":    serializers.SerializeEmails(emails),
	})
}
