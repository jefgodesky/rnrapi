package models

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	APIKey   string `json:"apikey"`
	Active   bool   `json:"active"`
}

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func HashAPIKey(apiKey string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckAPIKey(providedKey, storedHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedKey))
}
