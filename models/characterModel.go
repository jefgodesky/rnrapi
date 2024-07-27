package models

import (
	"gorm.io/gorm"
	"strings"
)

type Character struct {
	ID string `gorm:"primaryKey" json:"id"`
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Str         int    `json:"strength"`
	Dex         int    `json:"dexterity"`
	Con         int    `json:"constitution"`
	Int         int    `json:"intelligence"`
	Wis         int    `json:"wisdom"`
	Cha         int    `json:"charisma"`
	Notes       string `json:"notes"`
	PC          bool   `json:"pc"`
	Public      bool   `json:"public"`
	PlayerID    uint   `json:"player_id"`
	Player      User   `gorm:"foreignKey:PlayerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"player"`
}

func (char *Character) BeforeCreate(tx *gorm.DB) (err error) {
	return UniqueIDBeforeCreate(tx, &Character{}, &char.ID)
}

func IsValidAbility(ability string) bool {
	physical := "Strength Dexterity Constitution"
	mental := "Intelligence Wisdom Charisma"
	resistances := "Fortitude Reflexes Will"
	validAbilities := physical + " " + mental + " " + resistances
	return strings.Contains(validAbilities, ability)
}

func GetFortitude(char Character) int {
	return max(char.Str, char.Con)
}

func GetReflexes(char Character) int {
	return max(char.Dex, char.Int)
}

func GetWill(char Character) int {
	return max(char.Wis, char.Cha)
}

func GetAbility(char Character, ability string) int {
	switch ability {
	case "Strength":
		return char.Str
	case "Constitution":
		return char.Con
	case "Dexterity":
		return char.Dex
	case "Intelligence":
		return char.Int
	case "Wisdom":
		return char.Wis
	case "Charisma":
		return char.Cha
	case "Fortitude":
		return GetFortitude(char)
	case "Reflexes":
		return GetReflexes(char)
	case "Will":
		return GetWill(char)
	default:
		return 0
	}
}
