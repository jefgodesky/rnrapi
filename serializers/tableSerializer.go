package serializers

import (
	"github.com/jefgodesky/rnrapi/models"
)

type SerializedRow struct {
	Min     *int    `json:"min,omitempty"`
	Max     *int    `json:"max,omitempty"`
	Text    string  `json:"text"`
	Formula *string `json:"formula,omitempty"`
}

type SerializedTable struct {
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	DiceLabel   string          `json:"dice-label"`
	Formula     string          `json:"formula"`
	Ability     *string         `json:"ability,omitempty"`
	Cumulative  bool            `json:"cumulative"`
	Rows        []SerializedRow `json:"rows"`
	Public      bool            `json:"public"`
	Author      UserStub        `json:"author"`
}

type TableStub struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Author      UserStub `json:"author"`
}

func SerializeTableRow(row models.TableRow) SerializedRow {
	return SerializedRow{
		Min:     row.Min,
		Max:     row.Max,
		Text:    row.Text,
		Formula: row.Formula,
	}
}

func SerializeTable(table models.Table) SerializedTable {
	var rows = make([]SerializedRow, len(table.Rows))
	for i, row := range table.Rows {
		rows[i] = SerializeTableRow(row)
	}

	serialized := SerializedTable{
		Name:        table.Name,
		Slug:        table.Slug,
		Description: table.Description,
		DiceLabel:   table.DiceLabel,
		Formula:     table.Formula,
		Cumulative:  table.Cumulative,
		Rows:        rows,
		Public:      table.Public,
		Author:      StubUser(table.Author),
	}

	if table.Ability != nil && models.IsValidAbility(*table.Ability) {
		serialized.Ability = table.Ability
	}

	return serialized
}

func StubTable(table models.Table) TableStub {
	return TableStub{
		Name:        table.Name,
		Path:        "/tables/" + table.Slug,
		Description: table.Description,
		Author:      StubUser(table.Author),
	}
}

func SerializeTables(tables []models.Table) []TableStub {
	stubs := make([]TableStub, 0)
	for _, table := range tables {
		stubs = append(stubs, StubTable(table))
	}
	return stubs
}
