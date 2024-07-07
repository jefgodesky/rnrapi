package serializers

import (
	"encoding/json"
	"github.com/jefgodesky/rnrapi/models"
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
	Notes       []string              `json:"notes"`
	Public      bool                  `json:"public"`
	Player      string                `json:"player"`
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

	var notes []string
	if err := json.Unmarshal(char.Notes, &notes); err != nil {
		notes = []string{}
	}

	return SerializedCharacter{
		ID:          char.ID,
		Name:        char.Name,
		Description: char.Description,
		Abilities:   abilities,
		Resistances: resistances,
		Notes:       notes,
		Public:      char.Public,
		Player:      char.Player.Username,
	}
}

func SerializeCharacters(chars []models.Character) []SerializedCharacter {
	serializedChars := make([]SerializedCharacter, 0)
	for _, char := range chars {
		serializedChars = append(serializedChars, SerializeCharacter(char))
	}
	return serializedChars
}
