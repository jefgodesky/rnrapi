package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(plaintext string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckHash(plaintext, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
}
