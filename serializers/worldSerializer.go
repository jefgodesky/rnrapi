package serializers

import (
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedWorld struct {
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	Public      bool           `json:"public"`
	Creators    []string       `json:"creators"`
	Species     []SpeciesStub  `json:"species"`
	Societies   []SocietyStub  `json:"societies"`
	Campaigns   []CampaignStub `json:"campaigns"`
}

type WorldStub struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Description string     `json:"description"`
	Creators    []UserStub `json:"creators"`
}

func SerializeWorld(world models.World) SerializedWorld {
	creators := UsersToUsernames(world.Creators)
	var campaigns []models.Campaign
	initializers.DB.Where("world_id = ?", world.ID).Preload("GMs").Find(&campaigns)

	serializedCampaigns := make([]CampaignStub, len(campaigns))
	for i, campaign := range campaigns {
		serializedCampaigns[i] = StubCampaignWithWorld(campaign, world.Slug)
	}

	var species []models.Species
	initializers.DB.Where("world_id = ?", world.ID).Find(&species)

	serializedSpecies := make([]SpeciesStub, len(species))
	for i, sp := range species {
		serializedSpecies[i] = StubSpeciesWithWorld(sp, world.Slug)
	}

	var societies []models.Society
	initializers.DB.Where("world_id = ?", world.ID).Find(&societies)

	serializedSocieties := make([]SocietyStub, len(societies))
	for i, society := range societies {
		serializedSocieties[i] = StubSocietyWithWorld(society, world.Slug)
	}

	return SerializedWorld{
		Name:        world.Name,
		Slug:        world.Slug,
		Description: world.Description,
		Public:      world.Public,
		Creators:    creators,
		Species:     serializedSpecies,
		Societies:   serializedSocieties,
		Campaigns:   serializedCampaigns,
	}
}

func StubWorld(world models.World) WorldStub {
	serialized := SerializeWorld(world)
	creators := make([]UserStub, len(world.Creators))
	for i, creator := range world.Creators {
		creators[i] = StubUser(creator)
	}

	return WorldStub{
		Name:        serialized.Name,
		Path:        "/worlds/" + world.Slug,
		Description: world.Description,
		Creators:    creators,
	}
}

func SerializeWorlds(worlds []models.World) []WorldStub {
	stubs := make([]WorldStub, 0)
	for _, world := range worlds {
		stubs = append(stubs, StubWorld(world))
	}
	return stubs
}
