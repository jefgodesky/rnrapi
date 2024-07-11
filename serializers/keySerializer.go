package serializers

import (
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedKey struct {
	ID        uint    `json:"id"`
	Label     string  `json:"label"`
	Ephemeral bool    `json:"ephemeral"`
	Value     *string `json:"value"`
}

func SerializeKey(key models.Key, val *string) SerializedKey {
	serialized := SerializedKey{
		ID:        key.ID,
		Label:     key.Label,
		Ephemeral: key.Ephemeral,
	}

	if val != nil {
		serialized.Value = val
	}

	return serialized
}

func SerializeKeys(keys []models.Key) []SerializedKey {
	serialized := make([]SerializedKey, 0)
	for _, key := range keys {
		serialized = append(serialized, SerializeKey(key, nil))
	}
	return serialized
}
