package models

import (
	"errors"
	"fmt"
	"github.com/jefgodesky/rnrapi/enums"
	"gorm.io/gorm"
)

type Society struct {
	gorm.Model
	Slug        string            `gorm:"unique" json:"slug"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Favored     enums.AbilityPair `gorm:"type:string" json:"favored"`
	Languages   string            `json:"languages"`
	Public      bool              `json:"public"`
	WorldID     uint              `json:"world_id"`
	World       World             `gorm:"foreignKey:WorldID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"world"`
}

func (society *Society) BeforeSave(tx *gorm.DB) (err error) {
	for _, favored := range society.Favored {
		if !favored.IsValid() {
			return errors.New(fmt.Sprintf("invalid favored ability '%s'", favored))
		}
	}

	if society.Favored[0] == society.Favored[1] {
		return errors.New("favored abilities must be distinct")
	}

	return
}

func (Society) TableName() string {
	return "societies"
}
