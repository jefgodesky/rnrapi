package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GeneratePassword(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	password := base64.RawURLEncoding.EncodeToString(randomBytes)
	return password[:length], nil
}
