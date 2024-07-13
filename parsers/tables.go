package parsers

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jefgodesky/rnrapi/helpers"
	"github.com/jefgodesky/rnrapi/models"
	"gorm.io/gorm"
)

func BodyToTable(c *gin.Context) *models.Table {
	type TableRowBody struct {
		ID      uint    `json:"id"`
		Min     *int    `json:"min"`
		Max     *int    `json:"max"`
		Text    string  `json:"text"`
		Formula *string `json:"formula"`
	}

	var body struct {
		Name        string         `json:"name"`
		Slug        string         `json:"slug"`
		Description string         `json:"description"`
		DiceLabel   string         `json:"dice-label"`
		Formula     string         `json:"formula"`
		Ability     *string        `json:"ability,omitempty"`
		Cumulative  *bool          `json:"cumulative"`
		Rows        []TableRowBody `json:"rows"`
		Public      *bool          `json:"public"`
		Author      *string        `json:"author"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	tableSlug := body.Slug
	if tableSlug == "" {
		tableSlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	isCumulative := false
	if body.Cumulative != nil {
		isCumulative = *body.Cumulative
	}

	author := helpers.GetUserFromContext(c, true)
	if body.Author != nil {
		author = helpers.GetUser(c, *body.Author, false)
	}

	var rows = make([]models.TableRow, 0)
	for _, row := range body.Rows {
		rows = append(rows, models.TableRow{
			Model:   gorm.Model{ID: row.ID},
			Min:     row.Min,
			Max:     row.Max,
			Text:    row.Text,
			Formula: row.Formula,
		})
	}

	table := models.Table{
		Name:        body.Name,
		Slug:        tableSlug,
		Description: body.Description,
		DiceLabel:   body.DiceLabel,
		Formula:     body.Formula,
		Cumulative:  isCumulative,
		Rows:        rows,
		Public:      isPublic,
		Author:      *author,
	}

	if body.Ability != nil && models.IsValidAbility(*body.Ability) {
		table.Ability = body.Ability
	}

	return &table
}
