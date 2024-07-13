package parsers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
)

func UsernamesToUsers(usernames []string) []models.User {
	var users []models.User
	for _, username := range usernames {
		var user models.User
		result := initializers.DB.
			Where("username = ? AND active = ?", username, true).
			First(&user)
		if result.Error == nil {
			users = append(users, user)
		}
	}
	return users
}

func UsernamesToUsersWithDefault(usernames []string, defaultUser models.User) []models.User {
	users := UsernamesToUsers(usernames)
	if len(users) == 0 {
		users = append(users, defaultUser)
	}
	return users
}

func BodyToUserFields(c *gin.Context) (string, string, string, string) {
	var body struct {
		Username string
		Password *string
		Name     string
		Bio      string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return "", "", "", ""
	}

	password := *body.Password
	if body.Password == nil {
		password = ""
	}

	return body.Username, password, body.Name, body.Bio
}
