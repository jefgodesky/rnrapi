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

type CampaignStub struct {
	Name string `json:"name"`
	Path string `json:"path"`
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

func StubCampaignWithWorld(campaign models.Campaign, world string) CampaignStub {
	serialized := SerializeCampaign(campaign)
	return CampaignStub{
		Name: serialized.Name,
		Path: "/campaigns/" + world + "/" + campaign.Slug,
	}
}

func StubCampaign(campaign models.Campaign) CampaignStub {
	return StubCampaignWithWorld(campaign, campaign.World.Slug)
}

func SerializeCampaigns(campaigns []models.Campaign) []CampaignStub {
	stubs := make([]CampaignStub, 0)
	for _, campaign := range campaigns {
		stubs = append(stubs, StubCampaign(campaign))
	}
	return stubs
}
