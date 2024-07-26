package serializers

import (
	"github.com/jefgodesky/rnrapi/models"
	"sort"
)

type SerializedLevel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SerializedScale struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Levels      []SerializedLevel `json:"levels"`
	Author      UserStub          `json:"author"`
}

type ScaleStub struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Author      UserStub `json:"author"`
}

func SerializeLevel(level models.Level) SerializedLevel {
	return SerializedLevel{
		Name:        level.Name,
		Description: level.Description,
	}
}

func SerializeScale(scale models.Scale) SerializedScale {
	sortedLevels := make([]models.Level, len(scale.Levels))
	copy(sortedLevels, scale.Levels)
	sort.Slice(sortedLevels, func(i, j int) bool {
		return sortedLevels[i].Order < sortedLevels[j].Order
	})

	var levels = make([]SerializedLevel, len(sortedLevels))
	for i, level := range sortedLevels {
		levels[i] = SerializeLevel(level)
	}

	return SerializedScale{
		Name:        scale.Name,
		Slug:        scale.Slug,
		Description: scale.Description,
		Levels:      levels,
		Author:      StubUser(scale.Author),
	}
}

func StubScale(scale models.Scale) ScaleStub {
	return ScaleStub{
		Name:        scale.Name,
		Path:        "/scales/" + scale.Slug,
		Description: scale.Description,
		Author:      StubUser(scale.Author),
	}
}

func SerializeScales(scales []models.Scale) []ScaleStub {
	stubs := make([]ScaleStub, 0)
	for _, scale := range scales {
		stubs = append(stubs, StubScale(scale))
	}
	return stubs
}
