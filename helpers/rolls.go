package helpers

import (
	"github.com/jefgodesky/rnrapi/models"
	"github.com/justinian/dice"
	"regexp"
	"strconv"
	"strings"
)

func LogRoll(roll *models.Roll, res dice.RollResult) {
	prefix := models.RollLogSeparator
	if roll.Log == "" {
		prefix = ""
	}
	roll.Log += prefix + res.String()
}

func AddToResults(roll *models.Roll, results []string) {
	prefix := models.RollResultSeparator
	if roll.Results == "" {
		prefix = ""
	}
	roll.Results += prefix + strings.Join(results, models.RollResultSeparator)
}

func evaluateFormulaDice(formula string, roll *models.Roll) string {
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
		LogRoll(roll, res)
		formula = strings.Replace(formula, match[0], strconv.Itoa(res.Int()), 1)
	}

	return formula
}

func EvaluateFormula(formula string, roll *models.Roll) {
	formula = evaluateFormulaDice(formula, roll)
	AddToResults(roll, []string{formula})
}

func ProcessRow(row models.TableRow, roll *models.Roll) {
	if row.Formula != nil {
		EvaluateFormula(*row.Formula, roll)
	} else {
		AddToResults(roll, []string{row.Text})
	}
}

func CheckTable(table models.Table, number int, roll *models.Roll) {
	for _, row := range table.Rows {
		moreThanMin := row.Min == nil || number >= *row.Min
		lessThanMax := row.Max == nil || number <= *row.Max
		if (table.Cumulative && moreThanMin) || (moreThanMin && lessThanMax) {
			ProcessRow(row, roll)
		}
	}
}

func RollOnTable(table models.Table, roll *models.Roll, modifier int) {
	modifierStr := "+" + strconv.Itoa(modifier)
	res, _, _ := dice.Roll(table.Formula + modifierStr)
	LogRoll(roll, res)
	CheckTable(table, res.Int(), roll)
}
