package models

import (
	"errors"
	"fmt"
	"github.com/jefgodesky/rnrapi/enums"
	"gorm.io/gorm"
)

type Species struct {
	gorm.Model
	Slug        string            `gorm:"unique" json:"slug"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Affinities  enums.AbilityPair `gorm:"type:string" json:"affinities"`
	Aversion    enums.Ability     `gorm:"type:string" json:"aversion"`
	Public      bool              `json:"public"`
	WorldID     uint              `json:"world_id"`
	World       World             `gorm:"foreignKey:WorldID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"world"`
}

func (species *Species) BeforeSave(tx *gorm.DB) (err error) {
	for _, affinity := range species.Affinities {
		if !affinity.IsValid() {
			return errors.New(fmt.Sprintf("invalid affinity '%s'", affinity))
		}
	}

	if species.Affinities[0] == species.Affinities[1] {
		return errors.New("affinities must be distinct")
	}

	if !species.Aversion.IsValid() {
		return errors.New("invalid ability for Aversion")
	}

	return
}
