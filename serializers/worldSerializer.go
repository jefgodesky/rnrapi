package serializers

import (
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedWorld struct {
	Name      string                       `json:"name"`
	Slug      string                       `json:"slug"`
	Public    bool                         `json:"public"`
	Creators  []string                     `json:"creators"`
	Species   []SerializedSpeciesSansWorld `json:"species"`
	Societies []SerializedSocietySansWorld `json:"societies"`
	Campaigns []CampaignStub               `json:"campaigns"`
}

type WorldStub struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func SerializeWorld(world models.World) SerializedWorld {
	creators := make([]string, 0)
	for _, creator := range world.Creators {
		creators = append(creators, creator.Username)
	}

	var campaigns []models.Campaign
	initializers.DB.Where("world_id = ?", world.ID).Preload("GMs").Find(&campaigns)

	serializedCampaigns := make([]CampaignStub, len(campaigns))
	for i, campaign := range campaigns {
		serializedCampaigns[i] = StubCampaign(campaign)
	}

	var species []models.Species
	initializers.DB.Where("world_id = ?", world.ID).Find(&species)

	serializedSpecies := make([]SerializedSpeciesSansWorld, len(species))
	for i, sp := range species {
		serializedSpecies[i] = SerializeSpeciesSansWorld(sp)
	}

	var societies []models.Society
	initializers.DB.Where("world_id = ?", world.ID).Find(&societies)

	serializedSocieties := make([]SerializedSocietySansWorld, len(societies))
	for i, society := range societies {
		serializedSocieties[i] = SerializeSocietySansWorld(society)
	}

	return SerializedWorld{
		Name:      world.Name,
		Slug:      world.Slug,
		Public:    world.Public,
		Creators:  creators,
		Species:   serializedSpecies,
		Societies: serializedSocieties,
		Campaigns: serializedCampaigns,
	}
}

func StubWorld(world models.World) WorldStub {
	serialized := SerializeWorld(world)
	return WorldStub{
		Name: serialized.Name,
		Path: "/worlds/" + world.Slug,
	}
}

func SerializeWorlds(worlds []models.World) []WorldStub {
	stubs := make([]WorldStub, 0)
	for _, world := range worlds {
		stubs = append(stubs, StubWorld(world))
	}
	return stubs
}
