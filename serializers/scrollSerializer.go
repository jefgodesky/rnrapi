package serializers

import (
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedScroll struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Seals       uint                `json:"seals"`
	Writers     []UserStub          `json:"writers"`
	Readers     []UserStub          `json:"readers"`
	Public      bool                `json:"public"`
	Campaign    *SerializedCampaign `json:"campaign,omitempty"`
}

func SerializeScroll(scroll models.Scroll) SerializedScroll {
	serialized := SerializedScroll{
		ID:          scroll.ID,
		Name:        scroll.Name,
		Description: scroll.Description,
		Seals:       scroll.Seals,
		Writers:     SerializeUsers(scroll.Writers),
		Readers:     SerializeUsers(scroll.Readers),
		Public:      scroll.Public,
	}

	if scroll.Campaign != nil {
		ptr := SerializeCampaign(*scroll.Campaign)
		serialized.Campaign = &ptr
	}

	return serialized
}

func SerializeScrolls(scrolls []models.Scroll) []SerializedScroll {
	stubs := make([]SerializedScroll, 0)
	for _, scroll := range scrolls {
		stubs = append(stubs, SerializeScroll(scroll))
	}
	return stubs
}
