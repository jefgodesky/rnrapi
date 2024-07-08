package serializers

import (
	"github.com/jefgodesky/rnrapi/models"
	"strings"
	"time"
)

type SerializedRoll struct {
	ID        string         `json:"id"`
	Note      *string        `json:"note,omitempty"`
	Time      time.Time      `json:"time"`
	Table     TableStub      `json:"table"`
	Roller    *UserStub      `json:"roller,omitempty"`
	Character *CharacterStub `json:"character,omitempty"`
	Ability   *string        `json:"ability,omitempty"`
	Campaign  *CampaignStub  `json:"campaign,omitempty"`
	Modifier  int            `json:"modifier"`
	Log       []string       `json:"log"`
	Result    []string       `json:"results"`
}

func SerializeRoll(roll models.Roll) SerializedRoll {
	serialized := SerializedRoll{
		ID:       roll.ID,
		Time:     roll.CreatedAt,
		Table:    StubTable(roll.Table),
		Modifier: roll.Modifier,
		Log:      strings.Split(roll.Log, models.RollLogSeparator),
		Result:   strings.Split(roll.Results, models.RollResultSeparator),
	}

	if roll.Note != nil {
		serialized.Note = roll.Note
	}

	if roll.Roller != nil {
		ptr := StubUser(*roll.Roller)
		serialized.Roller = &ptr
	}

	if roll.Character != nil {
		ptr := StubCharacter(*roll.Character)
		serialized.Character = &ptr
	}

	if roll.Ability != nil {
		serialized.Ability = roll.Ability
	}

	if roll.Campaign != nil {
		ptr := StubCampaign(*roll.Campaign)
		serialized.Campaign = &ptr
	}

	return serialized
}

func SerializeRolls(rolls []models.Roll) []SerializedRoll {
	serializedRolls := make([]SerializedRoll, 0)
	for _, roll := range rolls {
		serializedRolls = append(serializedRolls, SerializeRoll(roll))
	}
	return serializedRolls
}
