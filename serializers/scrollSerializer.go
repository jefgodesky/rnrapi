package serializers

import (
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedScroll struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Seals       uint                `json:"seals"`
	Writers     []string            `json:"writers"`
	Readers     []string            `json:"readers"`
	Public      bool                `json:"public"`
	Campaign    *SerializedCampaign `json:"campaign,omitempty"`
}

type ScrollStub struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func SerializeScroll(scroll models.Scroll) SerializedScroll {
	serialized := SerializedScroll{
		ID:          scroll.ID,
		Name:        scroll.Name,
		Description: scroll.Description,
		Seals:       scroll.Seals,
		Writers:     UsersToUsernames(scroll.Writers),
		Readers:     UsersToUsernames(scroll.Readers),
		Public:      scroll.Public,
	}

	if scroll.Campaign != nil {
		ptr := SerializeCampaign(*scroll.Campaign)
		serialized.Campaign = &ptr
	}

	return serialized
}

func StubScroll(scroll models.Scroll) ScrollStub {
	return ScrollStub{
		Name: scroll.Name,
		Path: "/scrolls/" + scroll.ID,
	}
}

func SerializeScrolls(scrolls []models.Scroll) []ScrollStub {
	stubs := make([]ScrollStub, 0)
	for _, scroll := range scrolls {
		stubs = append(stubs, StubScroll(scroll))
	}
	return stubs
}
