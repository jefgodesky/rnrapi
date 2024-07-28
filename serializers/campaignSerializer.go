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

func SerializeCampaign(campaign models.Campaign) SerializedCampaign {
	var gms []UserStub
	for _, gm := range campaign.GMs {
		gms = append(gms, StubUser(gm))
	}

	var pcs []CharacterStub
	for _, pc := range campaign.PCs {
		pcs = append(pcs, StubCharacter(pc))
	}

	world := StubWorld(campaign.World)

	return SerializedCampaign{
		Name:        campaign.Name,
		Slug:        campaign.Slug,
		Description: campaign.Description,
		GMs:         gms,
		PCs:         pcs,
		Public:      campaign.Public,
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
