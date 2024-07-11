package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"-"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Active   bool   `json:"active"`
}
