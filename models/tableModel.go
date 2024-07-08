package models

import (
	"gorm.io/gorm"
	"strings"
)

type TableRow struct {
	gorm.Model
	Min     *int    `json:"min"`
	Max     *int    `json:"max"`
	Text    string  `json:"text"`
	Formula *string `json:"formula"`
	TableID uint    `json:"table_id"`
}

type Table struct {
	gorm.Model
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	DiceLabel   string     `json:"dice_label"`
	Formula     string     `json:"formula"`
	Ability     *string    `json:"ability,omitempty"`
	Cumulative  bool       `json:"cumulative"`
	Rows        []TableRow `gorm:"foreignKey:TableID" json:"rows"`
	Public      bool       `json:"public"`
	AuthorID    uint       `json:"author_id"`
	Author      User       `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"author"`
}

func IsValidAbility(ability string) bool {
	physical := "Strength Dexterity Constitution"
	mental := "Intelligence Wisdom Charisma"
	resistances := "Fortitude Reflexes Will"
	validAbilities := physical + " " + mental + " " + resistances
	return strings.Contains(validAbilities, ability)
}
