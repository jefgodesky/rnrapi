package serializers

import (
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedUser struct {
	Username   string          `json:"username"`
	Name       string          `json:"name"`
	Bio        string          `json:"bio"`
	Characters []CharacterStub `json:"characters"`
	APIKey     *string         `json:"api_key,omitempty"`
	Active     bool            `json:"active"`
}

type UserStub struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Path     string `json:"path"`
}

func SerializeUser(user models.User, apiKey *string) SerializedUser {
	var pcs []models.Character
	initializers.DB.Where("player_id = ? AND pc = ? AND public = ?", user.ID, true, true).Find(&pcs)

	characters := make([]CharacterStub, len(pcs))
	for i, pc := range pcs {
		characters[i] = StubCharacter(pc)
	}

	serializedUser := SerializedUser{
		Username:   user.Username,
		Name:       user.Name,
		Bio:        user.Bio,
		Characters: characters,
		Active:     user.Active,
	}

	if apiKey != nil {
		serializedUser.APIKey = apiKey
	}

	return serializedUser
}

func StubUser(user models.User) UserStub {
	return UserStub{
		Username: user.Username,
		Name:     user.Name,
		Path:     "/users/" + user.Username,
	}
}

func SerializeUsers(users []models.User) []UserStub {
	stubs := make([]UserStub, 0)
	for _, user := range users {
		stubs = append(stubs, StubUser(user))
	}
	return stubs
}
