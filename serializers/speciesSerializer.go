package serializers

import (
	"github.com/jefgodesky/rnrapi/enums"
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedStage struct {
	Name       string `json:"name"`
	Procedures string `json:"procedures"`
	MinAge     *uint  `json:"min"`
	MaxAge     *uint  `json:"max"`
}

type SerializedSpecies struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Affinities  enums.AbilityPair `json:"affinities"`
	Aversion    enums.Ability     `json:"aversion"`
	Stages      []SerializedStage `json:"stages"`
	Public      bool              `json:"public"`
	World       WorldStub         `json:"world"`
}

type SpeciesStub struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func SerializeStage(stage models.Stage) SerializedStage {
	return SerializedStage{
		Name:       stage.Name,
		Procedures: stage.Procedures,
		MinAge:     stage.MinAge,
		MaxAge:     stage.MaxAge,
	}
}

func SerializeSpecies(species models.Species) SerializedSpecies {
	world := StubWorld(species.World)
	serializedStages := make([]SerializedStage, len(species.Stages))
	for i, stage := range species.Stages {
		serializedStages[i] = SerializeStage(stage)
	}

	return SerializedSpecies{
		Name:        species.Name,
		Slug:        species.Slug,
		Description: species.Description,
		Affinities:  species.Affinities,
		Aversion:    species.Aversion,
		Stages:      serializedStages,
		Public:      species.Public,
		World:       world,
	}
}

func StubSpeciesWithWorld(species models.Species, world string) SpeciesStub {
	return SpeciesStub{
		Name: species.Name,
		Path: "/species/" + world + "/" + species.Slug,
	}
}

func StubSpecies(species models.Species) SpeciesStub {
	return StubSpeciesWithWorld(species, species.World.Slug)
}

func SerializeSpp(species []models.Species) []SpeciesStub {
	stubs := make([]SpeciesStub, 0)
	for _, sp := range species {
		stubs = append(stubs, StubSpecies(sp))
	}
	return stubs
}
