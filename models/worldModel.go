package models

import (
	"gorm.io/gorm"
)

type World struct {
	gorm.Model
	Slug     string `gorm:"unique" json:"slug"`
	Name     string `gorm:"unique" json:"name"`
	Creators []User `gorm:"many2many:world_creators;" json:"creators"`
	Public   bool   `json:"active" json:"public"`
}
