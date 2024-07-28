package serializers

import (
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedCampaign struct {
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	GMs         []UserStub      `json:"gms"`
	PCs         []CharacterStub `json:"pcs"`
	Public      bool            `json:"public"`
	World       WorldStub       `json:"world"`
}

type SerializedCampaignSansWorld struct {
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	GMs         []UserStub      `json:"gms"`
	PCs         []CharacterStub `json:"pcs"`
	Public      bool            `json:"public"`
}

func SerializeCampaignSansWorld(campaign models.Campaign) SerializedCampaignSansWorld {
	var gms []UserStub
	for _, gm := range campaign.GMs {
		gms = append(gms, StubUser(gm))
	}

	var pcs []CharacterStub
	for _, pc := range campaign.PCs {
		pcs = append(pcs, StubCharacter(pc))
	}

	return SerializedCampaignSansWorld{
		Name:        campaign.Name,
		Slug:        campaign.Slug,
		Description: campaign.Description,
		GMs:         gms,
		PCs:         pcs,
		Public:      campaign.Public,
	}
}

func SerializeCampaign(campaign models.Campaign) SerializedCampaign {
	sans := SerializeCampaignSansWorld(campaign)
	world := StubWorld(campaign.World)
	return SerializedCampaign{
		Name:        sans.Name,
		Slug:        sans.Slug,
		Description: sans.Description,
		GMs:         sans.GMs,
		PCs:         sans.PCs,
		Public:      sans.Public,
		World:       world,
	}
}

func SerializeCampaigns(campaigns []models.Campaign) []SerializedCampaign {
	stubs := make([]SerializedCampaign, 0)
	for _, campaign := range campaigns {
		stubs = append(stubs, SerializeCampaign(campaign))
	}
	return stubs
}
