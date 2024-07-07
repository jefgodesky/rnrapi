package serializers

import (
	"encoding/json"
	"github.com/jefgodesky/rnrapi/enums"
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedSpecies struct {
	Name        string                     `json:"name"`
	Slug        string                     `json:"slug"`
	Description string                     `json:"description"`
	Affinities  enums.AbilityPair          `json:"affinities"`
	Aversion    enums.Ability              `json:"aversion"`
	Stages      json.RawMessage            `json:"stages"`
	Public      bool                       `json:"public"`
	World       SerializedWorldSansSpecies `json:"world"`
}

type SerializedSpeciesSansWorld struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Affinities  enums.AbilityPair `json:"affinities"`
	Aversion    enums.Ability     `json:"aversion"`
	Stages      json.RawMessage   `json:"stages"`
	Public      bool              `json:"public"`
}

func SerializeSpecies(species models.Species) SerializedSpecies {
	world := SerializeWorldSansSpecies(species.World)
	return SerializedSpecies{
		Name:        species.Name,
		Slug:        species.Slug,
		Description: species.Description,
		Affinities:  species.Affinities,
		Aversion:    species.Aversion,
		Stages:      species.Stages,
		Public:      species.Public,
		World:       world,
	}
}

func SerializeSpeciesSansWorld(species models.Species) SerializedSpeciesSansWorld {
	return SerializedSpeciesSansWorld{
		Name:        species.Name,
		Slug:        species.Slug,
		Description: species.Description,
		Affinities:  species.Affinities,
		Aversion:    species.Aversion,
		Stages:      species.Stages,
		Public:      species.Public,
	}
}

func SerializeSpp(species []models.Species) []SerializedSpecies {
	serializedSpp := make([]SerializedSpecies, 0)
	for _, sp := range species {
		serializedSpp = append(serializedSpp, SerializeSpecies(sp))
	}
	return serializedSpp
}
