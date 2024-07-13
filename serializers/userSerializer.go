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
	Active     bool            `json:"active"`
}

type UserStub struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Path     string `json:"path"`
}

func SerializeUser(user models.User) SerializedUser {
	var pcs []models.Character
	initializers.DB.Where("player_id = ? AND pc = ? AND public = ?", user.ID, true, true).Find(&pcs)

	characters := make([]CharacterStub, len(pcs))
	for i, pc := range pcs {
		characters[i] = StubCharacter(pc)
	}

	return SerializedUser{
		Username:   user.Username,
		Name:       user.Name,
		Bio:        user.Bio,
		Characters: characters,
		Active:     user.Active,
	}
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

func UsersToUsernames(users []models.User) []string {
	var usernames = make([]string, len(users))
	for i, user := range users {
		usernames[i] = user.Username
	}
	return usernames
}
