package serializers

import "github.com/jefgodesky/rnrapi/models"

type SerializedCampaign struct {
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	GMs         []string  `json:"gms"`
	Public      bool      `json:"public"`
	World       WorldStub `json:"world"`
}

type SerializedCampaignSansWorld struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	GMs         []string `json:"gms"`
	Public      bool     `json:"public"`
}

func SerializeCampaign(campaign models.Campaign) SerializedCampaign {
	var gms []string
	for _, gm := range campaign.GMs {
		gms = append(gms, gm.Username)
	}

	world := StubWorld(campaign.World)

	return SerializedCampaign{
		Name:        campaign.Name,
		Slug:        campaign.Slug,
		Description: campaign.Description,
		GMs:         gms,
		Public:      campaign.Public,
		World:       world,
	}
}

func SerializeCampaignSansWorld(campaign models.Campaign) SerializedCampaignSansWorld {
	serialized := SerializeCampaign(campaign)
	return SerializedCampaignSansWorld{
		Name:        serialized.Name,
		Slug:        serialized.Slug,
		Description: serialized.Description,
		GMs:         serialized.GMs,
		Public:      serialized.Public,
	}
}

func SerializeCampaigns(campaigns []models.Campaign) []SerializedCampaign {
	serializedCampaigns := make([]SerializedCampaign, 0)
	for _, campaign := range campaigns {
		serializedCampaigns = append(serializedCampaigns, SerializeCampaign(campaign))
	}
	return serializedCampaigns
}
