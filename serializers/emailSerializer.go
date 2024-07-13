package serializers

import (
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedEmail struct {
	ID       uint   `json:"id"`
	Address  string `json:"address"`
	Verified bool   `json:"verified"`
}

func SerializeEmail(email models.Email) SerializedEmail {
	return SerializedEmail{
		ID:       email.ID,
		Address:  email.Address,
		Verified: email.Verified,
	}
}

func SerializeEmails(emails []models.Email) []SerializedEmail {
	serialized := make([]SerializedEmail, 0)
	for _, email := range emails {
		serialized = append(serialized, SerializeEmail(email))
	}
	return serialized
}
