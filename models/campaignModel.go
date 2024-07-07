package models

import (
	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GMs         []User `gorm:"many2many:campaign_gms;" json:"gms"`
	Public      bool   `json:"public"`
	WorldID     uint   `json:"world_id"`
	World       World  `gorm:"foreignKey:WorldID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"world"`
}
