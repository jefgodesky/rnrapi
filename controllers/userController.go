package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/jefgodesky/rnrapi/parsers"
	"github.com/jefgodesky/rnrapi/serializers"
	"gorm.io/gorm/clause"
	"strings"
)

func UserCreate(c *gin.Context) {
	username, password, name, bio := parsers.BodyToUserFields(c)
	if username == "" || password == "" || name == "" || bio == "" {
		return
	}

	hash, err := helpers.Hash(password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		c.Abort()
		return
	}

	user := models.User{Username: username, Password: hash, Name: name, Bio: bio, Active: true}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		usernameError := "duplicate key value violates unique constraint"
		if strings.Contains(result.Error.Error(), usernameError) {
			c.JSON(409, gin.H{"error": fmt.Sprintf("Username %s already exists", username)})
		} else {
			c.JSON(400, gin.H{"error": "Failed to create user"})
		}
		c.Abort()
		return
	}

	location := fmt.Sprintf("/v1/users/%s", user.Username)
	c.Header("Location", location)
	c.JSON(200, serializers.SerializeUser(user))
}

func UserIndex(c *gin.Context) {
	var users []models.User
	query := initializers.DB.Model(&models.User{}).
		Where("active = ?", true).
		Preload(clause.Associations)

	q := c.Query("q")
	if q != "" {
		likePattern := "%" + q + "%"
		query = query.Where("name LIKE ? OR username LIKE ?", likePattern, likePattern)
	}

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
	c.JSON(200, serializers.SerializeUser(*user))
}

func UserUpdate(c *gin.Context) {
	username, password, name, bio := parsers.BodyToUserFields(c)
	if username == "" || name == "" || bio == "" {
		return
	}

	user := helpers.GetUserFromContext(c, true)
	user.Username = username
	user.Name = name
	user.Bio = bio

	if len(password) > 0 {
		hash, err := helpers.Hash(password)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			c.Abort()
			return
		}

		user.Password = hash
	}

	if err := initializers.DB.Save(user).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to update user"})
		c.Abort()
		return
	}

	c.JSON(200, serializers.SerializeUser(*user))
}

func UserDestroy(c *gin.Context) {
	user := helpers.GetUserFromContext(c, true)

	if err := initializers.DB.Delete(user).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Failed to delete user"})
		return
	}

	c.JSON(200, serializers.SerializeUser(*user))
}

func SendPasswordRestEmail(email string, password string) error {
	from := "The Ruins & Revolutions Catalogue <catalogue@ruinsandrevolutions.com>"
	to := email
	subject := "Your password has been reset"

	lines := make([]string, 15)
	lines[0] = "Your password has been reset. Your new password is:"
	lines[1] = ""
	lines[2] = password
	lines[3] = ""
	lines[4] = "Visit https://ruinsandrevolutions.com/login to log in with your"
	lines[5] = "new password. Once logged in, you can change it by going to"
	lines[6] = "https://ruinsandrevolutions.com/profile"

	return helpers.SendEmail(from, to, subject, strings.Join(lines, "\n"))
}

func PasswordReset(c *gin.Context) {
	message := gin.H{"error": "Reset request received."}
	var body struct {
		Address string `json:"address"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return
	}

	var user models.User
	result := initializers.DB.
		Joins("JOIN emails ON emails.user_id = users.id").
		Where("emails.address = ? AND emails.verified = ?", body.Address, true).
		First(&user)

	if result.Error == nil {
		password, err := helpers.GeneratePassword(20)
		if err != nil {
			c.JSON(200, message)
			c.Abort()
			return
		}

		hash, err := helpers.Hash(password)
		if err != nil {
			c.JSON(200, message)
			c.Abort()
			return
		}

		user.Password = hash

		if err = initializers.DB.Save(user).Error; err != nil {
			c.JSON(200, message)
			c.Abort()
			return
		}

		err = SendPasswordRestEmail(body.Address, password)
		if err != nil {
			c.JSON(200, message)
			c.Abort()
			return
		}
	}

	c.JSON(200, message)
}
