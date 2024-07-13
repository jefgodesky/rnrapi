package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Email struct {
	gorm.Model
	Address  string `json:"address"`
	Code     string `json:"code"`
	Verified bool   `json:"verified"`
	UserID   uint   `json:"user_id"`
	User     User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
}

func SetVerificationCode(email *Email) {
	email.Code = uuid.New().String()
	email.Verified = false
}
