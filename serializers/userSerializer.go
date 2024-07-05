package serializers

import "github.com/jefgodesky/rnrapi/models"

type SerializedUser struct {
	Username string `json:"username"`
	Active   bool   `json:"active"`
}

func SerializeUser(user models.User) SerializedUser {
	return SerializedUser{
		Username: user.Username,
		Active:   user.Active,
	}
}
