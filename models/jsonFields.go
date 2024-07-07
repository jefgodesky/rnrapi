package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONField []byte

func (j *JSONField) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("scan source is not []byte")
	}
	return json.Unmarshal(bytes, j)
}

func (j *JSONField) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j JSONField) ToBytes() ([]byte, error) {
	return json.Marshal(j)
}
