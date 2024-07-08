package helpers

import (
	"fmt"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
	"github.com/justinian/dice"
	"gorm.io/gorm/clause"
	"regexp"
	"strconv"
	"strings"
)

func PrepareSubRoll(original *models.Roll, table *models.Table) models.Roll {
	sub := models.Roll{
		Note:     original.Note,
		Table:    *table,
		Modifier: 0,
	}

	if original.Roller != nil {
		sub.Roller = original.Roller
	}

	if original.Character != nil {
		sub.Character = original.Character
	}

	if original.Ability != nil {
		sub.Ability = original.Ability
	}

	if original.Campaign != nil {
		sub.Campaign = original.Campaign
	}

	return sub
}

func AddToLog(roll *models.Roll, logs []string) {
	prefix := models.RollLogSeparator
	if roll.Log == "" {
		prefix = ""
	}
	roll.Log += prefix + strings.Join(logs, models.RollLogSeparator)
}

func LogRoll(roll *models.Roll, res dice.RollResult) {
	AddToLog(roll, []string{res.String()})
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

func evaluateFormulaTables(formula string, roll *models.Roll) string {
	re := regexp.MustCompile(`\[([^\]]+)\]`)
	matches := re.FindAllStringSubmatch(formula, -1)
	if matches == nil {
		return formula
	}

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		tableSlug := match[1]
		var table models.Table
		result := initializers.DB.
			Preload(clause.Associations).
			Where("slug = ?", tableSlug).
			First(&table)
		if result.Error != nil {
			continue
		}

		subRoll := PrepareSubRoll(roll, &table)
		RollOnTable(table, &subRoll, 0)
		results := strings.Split(subRoll.Results, models.RollResultSeparator)
		resultsCS := strings.Join(results, ", ")
		formula = strings.Replace(formula, match[0], resultsCS, 1)
		AddToLog(roll, []string{subRoll.Log})
	}

	return formula
}

func EvaluateFormula(formula string, roll *models.Roll) {
	formula = evaluateFormulaDice(formula, roll)
	formula = evaluateFormulaTables(formula, roll)
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
	logMsg := fmt.Sprintf("Rolling on %s [%s]", table.Name, table.Slug)
	AddToLog(roll, []string{logMsg})

	modifierStr := "+" + strconv.Itoa(modifier)
	res, _, _ := dice.Roll(table.Formula + modifierStr)
	LogRoll(roll, res)
	CheckTable(table, res.Int(), roll)
}
