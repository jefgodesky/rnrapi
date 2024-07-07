package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"math/rand"
)

type Notes []string

func (n *Notes) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("scan source is not []byte")
	}
	return json.Unmarshal(bytes, n)
}

func (n *Notes) Value() (driver.Value, error) {
	return json.Marshal(n)
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
	Notes       json.RawMessage `gorm:"type:json" json:"notes"`
	Public      bool            `json:"public"`
	PlayerID    uint            `json:"player_id"`
	Player      User            `gorm:"foreignKey:PlayerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"player"`
}

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateID() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(b)
}

func isIDUnique(tx *gorm.DB, id string) bool {
	var count int64
	tx.Model(&Character{}).Where("id = ?", id).Count(&count)
	return count == 0
}

func (char *Character) BeforeCreate(tx *gorm.DB) (err error) {
	for {
		char.ID = generateID()
		if isIDUnique(tx, char.ID) {
			break
		}
	}

	notesJSON, err := json.Marshal(char.Notes)
	if err != nil {
		return err
	}
	char.Notes = notesJSON

	return
}
