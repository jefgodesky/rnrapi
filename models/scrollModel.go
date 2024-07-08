package models

import (
	"gorm.io/gorm"
)

type Scroll struct {
	ID string `gorm:"primaryKey" json:"id"`
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Seals       uint      `json:"seals"`
	Writers     []User    `gorm:"many2many:scroll_writers;" json:"writers"`
	Readers     []User    `gorm:"many2many:scroll_readers;" json:"readers"`
	Public      bool      `json:"active" json:"public"`
	CampaignID  *uint     `json:"campaign_id"`
	Campaign    *Campaign `gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"campaign"`
}

func (scroll *Scroll) BeforeCreate(tx *gorm.DB) (err error) {
	return UniqueIDBeforeCreate(tx, &Scroll{}, &scroll.ID)
}
