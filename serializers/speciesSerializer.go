package serializers

import (
	"github.com/jefgodesky/rnrapi/enums"
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedSpecies struct {
	Name        string                       `json:"name"`
	Slug        string                       `json:"slug"`
	Description string                       `json:"description"`
	Affinities  enums.AbilityPair            `json:"affinities"`
	Aversion    enums.Ability                `json:"aversion"`
	Public      bool                         `json:"public"`
	World       SerializedWorldSansCampaigns `json:"world"`
}

type SerializedSpeciesSansWorld struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Affinities  enums.AbilityPair `json:"affinities"`
	Aversion    enums.Ability     `json:"aversion"`
	Public      bool              `json:"public"`
}

func SerializeSpecies(species models.Species) SerializedSpecies {
	world := SerializeWorldSansCampaigns(species.World)

	return SerializedSpecies{
		Name:        species.Name,
		Slug:        species.Slug,
		Description: species.Description,
		Affinities:  species.Affinities,
		Aversion:    species.Aversion,
		Public:      world.Public,
		World:       world,
	}
}

func SerializeSpp(species []models.Species) []SerializedSpecies {
	serializedSpp := make([]SerializedSpecies, 0)
	for _, sp := range species {
		serializedSpp = append(serializedSpp, SerializeSpecies(sp))
	}
	return serializedSpp
}
