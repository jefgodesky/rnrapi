package models

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Key struct {
	gorm.Model
	Label     string `json:"label"`
	Token     string `gorm:"uniqueIndex" json:"-"`
	Secret    string `json:"-"`
	Ephemeral bool   `json:"ephemeral"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
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
