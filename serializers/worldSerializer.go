package serializers

import (
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedWorld struct {
	Name      string                        `json:"name"`
	Slug      string                        `json:"slug"`
	Public    bool                          `json:"public"`
	Creators  []string                      `json:"creators"`
	Species   []SerializedSpeciesSansWorld  `json:"species"`
	Societies []SerializedSocietySansWorld  `json:"societies"`
	Campaigns []SerializedCampaignSansWorld `json:"campaigns"`
}

type SerializedWorldSansCampaigns struct {
	Name      string                       `json:"name"`
	Slug      string                       `json:"slug"`
	Public    bool                         `json:"public"`
	Creators  []string                     `json:"creators"`
	Species   []SerializedSpeciesSansWorld `json:"species"`
	Societies []SerializedSocietySansWorld `json:"societies"`
}

type SerializedWorldSansSpecies struct {
	Name      string                        `json:"name"`
	Slug      string                        `json:"slug"`
	Public    bool                          `json:"public"`
	Creators  []string                      `json:"creators"`
	Societies []SerializedSocietySansWorld  `json:"societies"`
	Campaigns []SerializedCampaignSansWorld `json:"campaigns"`
}

type SerializedWorldSansSocieties struct {
	Name      string                        `json:"name"`
	Slug      string                        `json:"slug"`
	Public    bool                          `json:"public"`
	Creators  []string                      `json:"creators"`
	Species   []SerializedSpeciesSansWorld  `json:"species"`
	Campaigns []SerializedCampaignSansWorld `json:"campaigns"`
}

func SerializeWorld(world models.World) SerializedWorld {
	creators := make([]string, 0)
	for _, creator := range world.Creators {
		creators = append(creators, creator.Username)
	}

	var campaigns []models.Campaign
	initializers.DB.Where("world_id = ?", world.ID).Preload("GMs").Find(&campaigns)

	serializedCampaigns := make([]SerializedCampaignSansWorld, len(campaigns))
	for i, campaign := range campaigns {
		serializedCampaigns[i] = SerializeCampaignSansWorld(campaign)
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

func SerializeWorldSansCampaigns(world models.World) SerializedWorldSansCampaigns {
	serialized := SerializeWorld(world)
	return SerializedWorldSansCampaigns{
		Name:      serialized.Name,
		Slug:      serialized.Slug,
		Public:    serialized.Public,
		Creators:  serialized.Creators,
		Species:   serialized.Species,
		Societies: serialized.Societies,
	}
}

func SerializeWorldSansSpecies(world models.World) SerializedWorldSansSpecies {
	serialized := SerializeWorld(world)
	return SerializedWorldSansSpecies{
		Name:      serialized.Name,
		Slug:      serialized.Slug,
		Public:    serialized.Public,
		Creators:  serialized.Creators,
		Societies: serialized.Societies,
		Campaigns: serialized.Campaigns,
	}
}

func SerializeWorldSansSocieties(world models.World) SerializedWorldSansSocieties {
	serialized := SerializeWorld(world)
	return SerializedWorldSansSocieties{
		Name:      serialized.Name,
		Slug:      serialized.Slug,
		Public:    serialized.Public,
		Creators:  serialized.Creators,
		Species:   serialized.Species,
		Campaigns: serialized.Campaigns,
	}
}

func SerializeWorlds(worlds []models.World) []SerializedWorld {
	serializedWorlds := make([]SerializedWorld, 0)
	for _, world := range worlds {
		serializedWorlds = append(serializedWorlds, SerializeWorld(world))
	}
	return serializedWorlds
}
