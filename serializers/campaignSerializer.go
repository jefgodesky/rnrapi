package serializers

import "github.com/jefgodesky/rnrapi/models"

type SerializedCampaign struct {
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	GMs         []string        `json:"gms"`
	Public      bool            `json:"public"`
	World       SerializedWorld `json:"world"`
}

func SerializeCampaign(campaign models.Campaign) SerializedCampaign {
	var gms []string
	for _, gm := range campaign.GMs {
		gms = append(gms, gm.Username)
	}

	world := SerializeWorld(campaign.World)

	return SerializedCampaign{
		Name:        campaign.Name,
		Slug:        campaign.Slug,
		Description: campaign.Description,
		GMs:         gms,
		Public:      world.Public,
		World:       world,
	}
}
