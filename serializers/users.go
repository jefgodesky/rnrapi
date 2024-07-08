package serializers

import "github.com/jefgodesky/rnrapi/models"

func UsersToUsernames(users []models.User) []string {
	var usernames = make([]string, len(users))
	for i, user := range users {
		usernames[i] = user.Username
	}
	return usernames
}
