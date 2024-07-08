package helpers

import (
	"github.com/jefgodesky/rnrapi/models"
	"github.com/justinian/dice"
	"strconv"
	"strings"
)

func LogRoll(roll *models.Roll, res dice.RollResult, modifier string) {
	prefix := models.RollLogSeparator
	if roll.Log == "" {
		prefix = ""
	}
	roll.Log += prefix + res.String() + " Modifier: " + modifier
}

func AddToResults(roll *models.Roll, results []string) {
	prefix := models.RollResultSeparator
	if roll.Results == "" {
		prefix = ""
	}
	roll.Results += prefix + strings.Join(results, models.RollResultSeparator)
}

func ProcessRow(row models.TableRow) string {
	return row.Text
}

func CheckTable(table models.Table, number int) []string {
	results := make([]string, 0)
	for _, row := range table.Rows {
		moreThanMin := row.Min == nil || number >= *row.Min
		lessThanMax := row.Max == nil || number <= *row.Max
		if (table.Cumulative && moreThanMin) || (moreThanMin && lessThanMax) {
			results = append(results, ProcessRow(row))
		}
	}
	return results
}

func RollOnTable(table models.Table, roll *models.Roll, modifier int) {
	modifierStr := "+" + strconv.Itoa(modifier)
	res, _, _ := dice.Roll(table.Formula + modifierStr)
	LogRoll(roll, res, modifierStr)

	total := res.Int()
	results := CheckTable(table, total)
	AddToResults(roll, results)
}
