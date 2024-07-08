package models

import (
	"gorm.io/gorm"
)

type CharacterNote struct {
	gorm.Model
	Text        string    `json:"text"`
	CharacterID string    `json:"characterID"`
	Character   Character `gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE;" json:"-"`
}

type Character struct {
	ID string `gorm:"primaryKey" json:"id"`
	gorm.Model
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Str         int             `json:"strength"`
	Dex         int             `json:"dexterity"`
	Con         int             `json:"constitution"`
	Int         int             `json:"intelligence"`
	Wis         int             `json:"wisdom"`
	Cha         int             `json:"charisma"`
	Notes       []CharacterNote `json:"notes"`
	PC          bool            `json:"pc"`
	Public      bool            `json:"public"`
	PlayerID    uint            `json:"player_id"`
	Player      User            `gorm:"foreignKey:PlayerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"player"`
}

func (char *Character) BeforeCreate(tx *gorm.DB) (err error) {
	return UniqueIDBeforeCreate(tx, &Character{}, &char.ID)
}
