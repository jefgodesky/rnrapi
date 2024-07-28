package serializers

import (
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"gorm.io/gorm/clause"
)

type SerializedUser struct {
	Username   string               `json:"username"`
	Name       string               `json:"name"`
	Bio        string               `json:"bio"`
	Characters []CharacterStub      `json:"characters"`
	Campaigns  []SerializedCampaign `json:"campaigns"`
	Active     bool                 `json:"active"`
}

type UserStub struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Path     string `json:"path"`
}

func SerializeUser(user models.User) SerializedUser {
	var pcs []models.Character
	initializers.DB.
		Where("player_id = ? AND pc = ? AND public = ?", user.ID, true, true).
		Preload(clause.Associations).
		Find(&pcs)

	characters := make([]CharacterStub, len(pcs))
	for i, pc := range pcs {
		characters[i] = StubCharacter(pc)
	}

	var running []models.Campaign
	initializers.DB.
		Joins("JOIN campaign_gms ON campaign_gms.campaign_id = campaigns.id").
		Where("campaign_gms.user_id = ?", user.ID).
		Preload(clause.Associations).
		Find(&running)

	var playing []models.Campaign
	initializers.DB.
		Joins("JOIN campaign_pcs ON campaign_pcs.campaign_id = campaigns.id").
		Where("campaign_pcs.character_id IN (?)", initializers.DB.Select("id").Where("player_id = ?", user.ID).Table("characters")).
		Preload(clause.Associations).
		Find(&playing)

	campaignMap := make(map[uint]models.Campaign)
	for _, campaign := range running {
		campaignMap[campaign.ID] = campaign
	}

	for _, campaign := range playing {
		campaignMap[campaign.ID] = campaign
	}

	campaigns := make([]SerializedCampaign, 0, len(campaignMap))
	for _, campaign := range campaignMap {
		campaigns = append(campaigns, SerializeCampaign(campaign))
	}

	return SerializedUser{
		Username:   user.Username,
		Name:       user.Name,
		Bio:        user.Bio,
		Characters: characters,
		Campaigns:  campaigns,
		Active:     user.Active,
	}
}

func StubUser(user models.User) UserStub {
	return UserStub{
		Username: user.Username,
		Name:     user.Name,
		Path:     "/users/" + user.Username,
	}
}

func SerializeUsers(users []models.User) []UserStub {
	stubs := make([]UserStub, 0)
	for _, user := range users {
		stubs = append(stubs, StubUser(user))
	}
	return stubs
}
