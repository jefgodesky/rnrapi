package serializers

import (
	"github.com/jefgodesky/rnrapi/enums"
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedSociety struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Favored     enums.AbilityPair `json:"favored"`
	Languages   string            `json:"languages"`
	Public      bool              `json:"public"`
	World       WorldStub         `json:"world"`
}

type SerializedSocietySansWorld struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Favored     enums.AbilityPair `json:"favored"`
	Languages   string            `json:"languages"`
	Public      bool              `json:"public"`
}

func SerializeSociety(society models.Society) SerializedSociety {
	world := StubWorld(society.World)
	return SerializedSociety{
		Name:        society.Name,
		Slug:        society.Slug,
		Description: society.Description,
		Favored:     society.Favored,
		Languages:   society.Languages,
		Public:      society.Public,
		World:       world,
	}
}

func SerializeSocietySansWorld(society models.Society) SerializedSocietySansWorld {
	return SerializedSocietySansWorld{
		Name:        society.Name,
		Slug:        society.Slug,
		Description: society.Description,
		Favored:     society.Favored,
		Languages:   society.Languages,
		Public:      society.Public,
	}
}

func SerializeSocieties(societies []models.Society) []SerializedSociety {
	serializedSocieties := make([]SerializedSociety, 0)
	for _, society := range societies {
		serializedSocieties = append(serializedSocieties, SerializeSociety(society))
	}
	return serializedSocieties
}
