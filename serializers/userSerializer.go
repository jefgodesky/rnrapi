package serializers

import "github.com/jefgodesky/rnrapi/models"

type SerializedUser struct {
	Username string  `json:"username"`
	APIKey   *string `json:"api_key,omitempty"`
	Active   bool    `json:"active"`
}

func SerializeUser(user models.User, apiKey *string) SerializedUser {
	serializedUser := SerializedUser{
		Username: user.Username,
		Active:   user.Active,
	}

	if apiKey != nil {
		serializedUser.APIKey = apiKey
	}

	return serializedUser
}

func SerializeUsers(users []models.User) []SerializedUser {
	serializedUsers := make([]SerializedUser, 0)
	for _, user := range users {
		serializedUser := SerializeUser(user, nil)
		serializedUsers = append(serializedUsers, serializedUser)
	}
	return serializedUsers
}
