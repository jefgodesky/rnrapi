package helpers

import (
	"github.com/jefgodesky/rnrapi/models"
	"github.com/justinian/dice"
	"regexp"
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

func EvaluateFormula(formula string, roll *models.Roll) string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(formula, -1)
	if matches == nil {
		return formula
	}

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		res, _, _ := dice.Roll(match[1])
		LogRoll(roll, res, "+0")
		total := res.Int()
		formula = strings.Replace(formula, match[0], strconv.Itoa(total), 1)
	}

	return formula
}

func ProcessRow(row models.TableRow, roll *models.Roll) string {
	str := row.Text
	if row.Formula != nil {
		str = EvaluateFormula(*row.Formula, roll)
	}
	return str
}

func CheckTable(table models.Table, number int, roll *models.Roll) []string {
	results := make([]string, 0)
	for _, row := range table.Rows {
		moreThanMin := row.Min == nil || number >= *row.Min
		lessThanMax := row.Max == nil || number <= *row.Max
		if (table.Cumulative && moreThanMin) || (moreThanMin && lessThanMax) {
			results = append(results, ProcessRow(row, roll))
		}
	}
	return results
}

func RollOnTable(table models.Table, roll *models.Roll, modifier int) {
	modifierStr := "+" + strconv.Itoa(modifier)
	res, _, _ := dice.Roll(table.Formula + modifierStr)
	LogRoll(roll, res, modifierStr)

	total := res.Int()
	results := CheckTable(table, total, roll)
	AddToResults(roll, results)
}
