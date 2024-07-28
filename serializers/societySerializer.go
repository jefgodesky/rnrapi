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

type SocietyStub struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	World       WorldStub `json:"world"`
}

type SocietyStubSansWorld struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
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

func StubSocietyWithWorld(society models.Society, world string) SocietyStub {
	return SocietyStub{
		Name:        society.Name,
		Path:        "/societies/" + world + "/" + society.Slug,
		Description: society.Description,
		World:       StubWorld(society.World),
	}
}

func StubSociety(society models.Society) SocietyStub {
	return StubSocietyWithWorld(society, society.World.Slug)
}

func StubSocietySansWorld(society models.Society) SocietyStubSansWorld {
	return SocietyStubSansWorld{
		Name:        society.Name,
		Slug:        society.Slug,
		Description: society.Description,
	}
}

func SerializeSocieties(societies []models.Society) []SocietyStub {
	stubs := make([]SocietyStub, 0)
	for _, society := range societies {
		stubs = append(stubs, StubSociety(society))
	}
	return stubs
}
