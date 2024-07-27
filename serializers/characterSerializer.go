package serializers

import (
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"gorm.io/gorm/clause"
)

type SerializedAbilities struct {
	Str int `json:"strength"`
	Dex int `json:"dexterity"`
	Con int `json:"constitution"`
	Int int `json:"intelligence"`
	Wis int `json:"wisdom"`
	Cha int `json:"charisma"`
}

type SerializedResistances struct {
	Fort int `json:"fortitude"`
	Ref  int `json:"reflexes"`
	Will int `json:"will"`
}

type SerializedCharacter struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Abilities   SerializedAbilities   `json:"abilities"`
	Resistances SerializedResistances `json:"resistances"`
	Notes       string                `json:"notes"`
	PC          bool                  `json:"pc"`
	Campaigns   []CampaignStub        `json:"campaigns"`
	Public      bool                  `json:"public"`
	Player      string                `json:"player"`
}

type CharacterStub struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func SerializeCharacter(char models.Character) SerializedCharacter {
	abilities := SerializedAbilities{
		Str: char.Str,
		Dex: char.Dex,
		Con: char.Con,
		Int: char.Int,
		Wis: char.Wis,
		Cha: char.Cha,
	}

	resistances := SerializedResistances{
		Fort: max(char.Str, char.Con),
		Ref:  max(char.Dex, char.Int),
		Will: max(char.Wis, char.Cha),
	}

	var campaigns []models.Campaign
	err := initializers.DB.Joins("JOIN campaign_pcs ON campaign_pcs.campaign_id = campaigns.id").
		Preload(clause.Associations).
		Where("campaign_pcs.character_id = ?", char.ID).
		Find(&campaigns).Error
	if err != nil {
		panic("Could not find character campaigns")
	}

	var campaignStubs []CampaignStub
	for _, campaign := range campaigns {
		campaignStubs = append(campaignStubs, StubCampaign(campaign))
	}

	return SerializedCharacter{
		ID:          char.ID,
		Name:        char.Name,
		Description: char.Description,
		Abilities:   abilities,
		Resistances: resistances,
		Notes:       char.Notes,
		PC:          char.PC,
		Campaigns:   campaignStubs,
		Public:      char.Public,
		Player:      char.Player.Username,
	}
}

func StubCharacter(char models.Character) CharacterStub {
	return CharacterStub{
		Name: char.Name,
		Path: "/characters/" + char.ID,
	}
}

func SerializeCharacters(chars []models.Character) []CharacterStub {
	stubs := make([]CharacterStub, 0)
	for _, char := range chars {
		stubs = append(stubs, StubCharacter(char))
	}
	return stubs
}
