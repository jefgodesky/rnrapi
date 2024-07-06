package serializers

import "github.com/jefgodesky/rnrapi/models"

type SerializedWorld struct {
	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	Public   bool     `json:"public"`
	Creators []string `json:"creators"`
}

func SerializeWorld(world models.World) SerializedWorld {
	creators := []string{}
	for _, creator := range world.Creators {
		creators = append(creators, creator.Username)
	}

	return SerializedWorld{
		Name:     world.Name,
		Slug:     world.Slug,
		Public:   world.Public,
		Creators: creators,
	}
}
