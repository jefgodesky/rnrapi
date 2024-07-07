package helpers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jefgodesky/rnrapi/enums"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/models"
)

func UsernamesToUsers(usernames []string) []models.User {
	var users []models.User
	for _, username := range usernames {
		var user models.User
		result := initializers.DB.
			Where("username = ? AND active = ?", username, true).
			First(&user)
		if result.Error == nil {
			users = append(users, user)
		}
	}
	return users
}

func UsernamesToUsersWithDefault(usernames []string, defaultUser models.User) []models.User {
	users := UsernamesToUsers(usernames)
	if len(users) == 0 {
		users = append(users, defaultUser)
	}
	return users
}

func IdsToCharacters(ids []string) []models.Character {
	var characters []models.Character
	for _, id := range ids {
		var char models.Character
		result := initializers.DB.First(&char, "id = ?", id)
		if result.Error == nil {
			characters = append(characters, char)
		}
	}
	return characters
}

func BodyToUserFields(c *gin.Context) (string, string, string) {
	var body struct {
		Username string
		Name     string
		Bio      string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return "", "", ""
	}

	return body.Username, body.Name, body.Bio
}

func BodyToWorld(c *gin.Context) *models.World {
	var body struct {
		Name     string   `json:"name"`
		Slug     string   `json:"slug"`
		Public   *bool    `json:"public"`
		Creators []string `json:"creators"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return nil
	}

	worldSlug := body.Slug
	if worldSlug == "" {
		worldSlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	var creators []models.User
	for _, creator := range body.Creators {
		var user models.User
		result := initializers.DB.
			Where("username = ? AND active = ?", creator, true).
			First(&user)
		if result.Error == nil {
			creators = append(creators, user)
		}
	}

	if len(creators) == 0 {
		authUser := GetUserFromContext(c, true)
		creators = append(creators, *authUser)
	}

	world := models.World{
		Name:     body.Name,
		Slug:     worldSlug,
		Public:   isPublic,
		Creators: creators,
	}

	return &world
}

func BodyToCampaign(c *gin.Context) *models.Campaign {
	var body struct {
		Slug        string   `json:"slug"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		GMs         []string `json:"gms"`
		PCs         []string `json:"pcs"`
		Public      *bool    `json:"public"`
		World       string   `json:"world"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return nil
	}

	campaignSlug := body.Slug
	if campaignSlug == "" {
		campaignSlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	world := GetWorld(c, body.World)
	if world == nil {
		c.JSON(400, gin.H{"error": "World not found"})
		return nil
	}

	authUser := GetUserFromContext(c, true)
	gms := UsernamesToUsersWithDefault(body.GMs, *authUser)
	pcs := IdsToCharacters(body.PCs)

	campaign := models.Campaign{
		Slug:        campaignSlug,
		Name:        body.Name,
		Description: body.Description,
		GMs:         gms,
		PCs:         pcs,
		Public:      isPublic,
		WorldID:     world.ID,
		World:       *world,
	}

	return &campaign
}

func BodyToSpecies(c *gin.Context) *models.Species {
	var body struct {
		Slug        string           `json:"slug"`
		Name        string           `json:"name"`
		Description string           `json:"description"`
		Affinities  [2]enums.Ability `json:"affinities"`
		Aversion    enums.Ability    `json:"aversion"`
		Stages      json.RawMessage  `json:"stages"`
		Public      *bool            `json:"public"`
		World       string           `json:"world"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	speciesSlug := body.Slug
	if speciesSlug == "" {
		speciesSlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	world := GetWorld(c, body.World)
	if world == nil {
		c.JSON(400, gin.H{"error": "World not found"})
		c.Abort()
		return nil
	}

	species := models.Species{
		Slug:        speciesSlug,
		Name:        body.Name,
		Description: body.Description,
		Affinities:  enums.AbilityPair{body.Affinities[0], body.Affinities[1]},
		Aversion:    body.Aversion,
		Stages:      body.Stages,
		Public:      isPublic,
		WorldID:     world.ID,
		World:       *world,
	}

	return &species
}

func BodyToSociety(c *gin.Context) *models.Society {
	var body struct {
		Slug        string           `json:"slug"`
		Name        string           `json:"name"`
		Description string           `json:"description"`
		Favored     [2]enums.Ability `json:"favored"`
		Languages   string           `json:"languages"`
		Public      *bool            `json:"public"`
		World       string           `json:"world"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	societySlug := body.Slug
	if societySlug == "" {
		societySlug = slug.Make(body.Name)
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	world := GetWorld(c, body.World)
	if world == nil {
		c.JSON(400, gin.H{"error": "World not found"})
		c.Abort()
		return nil
	}

	society := models.Society{
		Slug:        societySlug,
		Name:        body.Name,
		Description: body.Description,
		Favored:     enums.AbilityPair{body.Favored[0], body.Favored[1]},
		Languages:   body.Languages,
		Public:      isPublic,
		WorldID:     world.ID,
		World:       *world,
	}

	return &society
}

func BodyToCharacter(c *gin.Context) *models.Character {
	type AbilitiesBody struct {
		Strength     int `json:"strength"`
		Dexterity    int `json:"dexterity"`
		Constitution int `json:"constitution"`
		Intelligence int `json:"intelligence"`
		Wisdom       int `json:"wisdom"`
		Charisma     int `json:"charisma"`
	}

	var body struct {
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Abilities   AbilitiesBody `json:"abilities"`
		Notes       []string      `json:"notes"`
		Public      *bool         `json:"public"`
		PC          bool          `json:"pc"`
		Player      string        `json:"player"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return nil
	}

	isPublic := true
	if body.Public != nil {
		isPublic = *body.Public
	}

	player := GetUser(c, body.Player)
	if player == nil {
		player = GetUserFromContext(c, true)
	}

	notesJSON, err := json.Marshal(body.Notes)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal notes"})
		return nil
	}

	char := models.Character{
		Name:        body.Name,
		Description: body.Description,
		Str:         body.Abilities.Strength,
		Dex:         body.Abilities.Dexterity,
		Con:         body.Abilities.Constitution,
		Int:         body.Abilities.Intelligence,
		Wis:         body.Abilities.Wisdom,
		Cha:         body.Abilities.Charisma,
		Notes:       notesJSON,
		PC:          body.PC,
		Public:      isPublic,
		PlayerID:    player.ID,
		Player:      *player,
	}

	return &char
}
