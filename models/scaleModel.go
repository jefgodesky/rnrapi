package models

import (
	"gorm.io/gorm"
)

type Level struct {
	gorm.Model
	Order       int    `json:"order"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ScaleID     uint   `json:"scale_id"`
}

type Scale struct {
	gorm.Model
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Public      bool    `json:"public"`
	Levels      []Level `gorm:"foreignKey:ScaleID" json:"levels"`
	AuthorID    uint    `json:"author_id"`
	Author      User    `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"author"`
}
