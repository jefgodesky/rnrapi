package models

import (
	"fmt"
	"github.com/justinian/dice"
	"gorm.io/gorm"
	"strconv"
)

const RollLogSeparator = "[ ~~~ ROLL LOG SEPARATOR ~~~]"
const RollResultSeparator = "[ ~~~ ROLL RESULT SEPARATOR ~~~]"

type Roll struct {
	ID string `gorm:"primaryKey" json:"id"`
	gorm.Model
	Note        *string    `json:"note"`
	TableID     uint       `json:"table_id"`
	Table       Table      `gorm:"foreignKey:TableID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"table"`
	RollerID    *uint      `json:"roller_id,omitempty"`
	Roller      *User      `gorm:"foreignKey:RollerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"roller,omitempty"`
	CharacterID *uint      `json:"character_id,omitempty"`
	Character   *Character `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"character,omitempty"`
	Ability     *string    `json:"ability,omitempty"`
	CampaignID  *uint      `json:"campaign_id,omitempty"`
	Campaign    *Campaign  `gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"campaign,omitempty"`
	Modifier    int        `json:"modifier"`
	Log         string     `json:"log"`
	Results     string     `json:"result"`
}

func (roll *Roll) BeforeCreate(tx *gorm.DB) (err error) {
	return UniqueIDBeforeCreate(tx, &Roll{}, &roll.ID)
}

func MakeRoll(instance Roll) {
	modifierStr := "+" + strconv.Itoa(instance.Modifier)
	res, _, _ := dice.Roll(instance.Table.Formula + modifierStr)
	total := res.Int()
	log := res.String() + " Modifier: " + modifierStr
	rolls := make([]string, 0)
	rolls = append(rolls, log)
	fmt.Println(total)
	fmt.Println(log)
}
