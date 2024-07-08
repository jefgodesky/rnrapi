package models

import (
	"gorm.io/gorm"
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
	Cumulative  bool       `json:"cumulative"`
	Rows        []TableRow `gorm:"foreignKey:TableID" json:"rows"`
	Public      bool       `json:"public"`
	AuthorID    uint       `json:"author_id"`
	Author      User       `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"author"`
}
