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
	Campaigns []SerializedCampaignSansWorld `json:"campaigns"`
}

type SerializedWorldSansCampaigns struct {
	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	Public   bool     `json:"public"`
	Creators []string `json:"creators"`
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

	return SerializedWorld{
		Name:      world.Name,
		Slug:      world.Slug,
		Public:    world.Public,
		Creators:  creators,
		Campaigns: serializedCampaigns,
	}
}

func SerializeWorldSansCampaigns(world models.World) SerializedWorldSansCampaigns {
	serialized := SerializeWorld(world)
	return SerializedWorldSansCampaigns{
		Name:     serialized.Name,
		Slug:     serialized.Slug,
		Public:   serialized.Public,
		Creators: serialized.Creators,
	}
}

func SerializeWorlds(worlds []models.World) []SerializedWorld {
	serializedWorlds := make([]SerializedWorld, 0)
	for _, world := range worlds {
		serializedWorlds = append(serializedWorlds, SerializeWorld(world))
	}
	return serializedWorlds
}
