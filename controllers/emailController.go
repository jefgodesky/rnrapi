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

	p := make([]string, 11)
	p[0] = "The only thing we will ever use your email address for is to "
	p[1] = "help you reset your password. If you lose your password or you "
	p[2] = "can’t remember it, we can create a new one for you and send it "
	p[3] = "to this email address. This works out great, so long as we can "
	p[4] = "be sure that when we send that new password to this email "
	p[5] = "address that we’re sending it to *you* and not some sinister "
	p[6] = "villain, like a shapeshifting doppelgänger, or any other, less "
	p[7] = "likely threat that may be lurking out there. That’s why we’d "
	p[8] = "like you to click the link below. That will verify that this "
	p[9] = "email reached you, so we know where to send your new password "
	p[10] = "if you should ever need it."
	paragraph := strings.Join(p, "")

	initializers.DB.Where("id = ?", email.UserID).First(&email.User)
	lines := make([]string, 5)
	lines[0] = fmt.Sprintf("Hello %s,", email.User.Name)
	lines[1] = ""
	lines[2] = paragraph
	lines[3] = ""
	lines[4] = fmt.Sprintf("https://ruinsandrevolutions.com/verify/%s-%d", email.Code, email.ID)

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

func EmailRetrieve(c *gin.Context) {
	email := helpers.GetEmailFromID(c)
	user := helpers.GetUserFromContext(c, true)
	if user == nil {
		return
	}

	if user.ID != email.UserID {
		c.JSON(403, gin.H{"error": "Forbidden"})
	}

	c.JSON(200, serializers.SerializeEmail(*email))
}

func EmailDestroy(c *gin.Context) {
	email := helpers.GetEmailFromID(c)
	user := helpers.GetUserFromContext(c, true)
	if user == nil {
		return
	}

	if user.ID != email.UserID {
		c.JSON(403, gin.H{"error": "Forbidden"})
	}

	if err := initializers.DB.Delete(email).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to delete key"})
		c.Abort()
		return
	}

	c.Status(204)
}

func EmailVerify(c *gin.Context) {
	code := parsers.BodyToVerification(c)
	if code == nil {
		return
	}

	email := helpers.GetEmailFromID(c)
	if email.Code != *code {
		c.JSON(400, gin.H{"error": "Invalid verification code"})
		c.Abort()
		return
	}

	email.Code = ""
	email.Verified = true

	if err := initializers.DB.Save(email).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update email record"})
		return
	}

	c.JSON(200, serializers.SerializeEmail(*email))
}
