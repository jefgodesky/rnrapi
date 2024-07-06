package models

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Token    string `gorm:"uniqueIndex" json:"-"`
	Secret   string `json:"-"`
	Active   bool   `json:"active"`
}

func GenerateAPIKey() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	secret := base64.StdEncoding.EncodeToString(bytes)
	token := uuid.New().String()
	return token, secret, nil
}

func HashAPIKey(apiKey string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckAPIKey(providedKey, storedHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedKey))
}
