package enums

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type Ability string

const (
	Strength     Ability = "strength"
	Dexterity    Ability = "dexterity"
	Constitution Ability = "constitution"
	Intelligence Ability = "intelligence"
	Wisdom       Ability = "wisdom"
	Charisma     Ability = "charisma"
)

func (ability *Ability) IsValid() bool {
	switch *ability {
	case Strength, Dexterity, Constitution, Intelligence, Wisdom, Charisma:
		return true
	}
	return false
}

func (ability *Ability) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid value for Ability")
	}
	*ability = Ability(str)
	return nil
}

func (ability *Ability) Value() (driver.Value, error) {
	return string(*ability), nil
}

type AbilityPair [2]Ability

func (pair *AbilityPair) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid value for AbilityPair")
	}

	parts := strings.Split(str, " ")
	if len(parts) != 2 {
		return errors.New("invalid length for AbilityPair")
	}

	(*pair)[0], (*pair)[1] = Ability(parts[0]), Ability(parts[1])
	return nil
}

func (pair AbilityPair) Value() (driver.Value, error) {
	return strings.Join([]string{string(pair[0]), string(pair[1])}, " "), nil
}
