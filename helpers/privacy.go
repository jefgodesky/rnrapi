package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/models"
	"reflect"
)

func IsWorldCreator(world *models.World, user *models.User) bool {
	if user == nil || world == nil {
		return false
	}

	for _, creator := range world.Creators {
		if creator.ID == user.ID {
			return true
		}
	}

	return false
}

func IsScrollUser(scroll *models.Scroll, user *models.User, role string) bool {
	if user == nil || scroll == nil {
		return false
	}

	list := scroll.Readers
	if role == "Writers" {
		list = scroll.Writers
	}

	for _, u := range list {
		if u.ID == user.ID {
			return true
		}
	}

	return false
}

func IsScrollReader(scroll *models.Scroll, user *models.User) bool {
	return IsScrollUser(scroll, user, "Readers")
}

func IsScrollWriter(scroll *models.Scroll, user *models.User) bool {
	return IsScrollUser(scroll, user, "Writers")
}

func HasWorldAccess(world *models.World, user *models.User) bool {
	if world.Public {
		return true
	}

	if IsWorldCreator(world, user) {
		return true
	}

	return false
}

func WorldCreatorOnly(c *gin.Context) *models.World {
	world := GetWorldFromSlug(c)
	if world == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	isCreator := IsWorldCreator(world, user)
	if !isCreator {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return world
}

func CharacterPlayerOnly(c *gin.Context) *models.Character {
	char := GetCharacterFromID(c)
	if char == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	if char.PlayerID != user.ID {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return char
}

func ScrollWriterOnly(c *gin.Context) *models.Scroll {
	scroll := GetScrollFromID(c)
	if scroll == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	isWriter := IsScrollWriter(scroll, user)
	if !isWriter {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return scroll
}

func TableAuthorOnly(c *gin.Context) *models.Table {
	table := GetTableFromSlug(c)
	if table == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	if table.AuthorID != user.ID {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return table
}

func ScaleAuthorOnly(c *gin.Context) *models.Scale {
	scale := GetScaleFromSlug(c)
	if scale == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	if scale.AuthorID != user.ID {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return scale
}

func filterWorldAccess(items interface{}, user *models.User) interface{} {
	itemsVal := reflect.ValueOf(items)
	if itemsVal.Kind() != reflect.Slice {
		panic("filterWorldAccess: items is not slice")
	}

	itemType := itemsVal.Type().Elem()
	filteredItems := reflect.MakeSlice(reflect.SliceOf(itemType), 0, itemsVal.Len())

	for i := 0; i < itemsVal.Len(); i++ {
		item := itemsVal.Index(i)
		worldField := item.FieldByName("World")
		if !worldField.IsValid() {
			panic("filterWorldAccess: No World found")
		}

		world, ok := worldField.Addr().Interface().(*models.World)
		if !ok {
			panic("filterWorldAccess: World is not of type *models.World")
		}

		if HasWorldAccess(world, user) {
			filteredItems = reflect.Append(filteredItems, item)
		}
	}

	return filteredItems.Interface()
}

func IsCampaignGM(campaign *models.Campaign, user *models.User) bool {
	if user == nil || campaign == nil {
		return false
	}

	for _, gm := range campaign.GMs {
		if gm.ID == user.ID {
			return true
		}
	}

	return false
}

func HasCampaignAccess(campaign *models.Campaign, user *models.User) bool {
	if !HasWorldAccess(&campaign.World, user) {
		return false
	}

	if campaign.Public {
		return true
	}

	if IsCampaignGM(campaign, user) {
		return true
	}

	return false
}

func CampaignGMOnly(c *gin.Context) *models.Campaign {
	campaign := GetCampaignFromSlug(c)
	if campaign == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	if !HasCampaignAccess(campaign, user) || !IsCampaignGM(campaign, user) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return campaign
}

func SpeciesCreatorOnly(c *gin.Context) *models.Species {
	species := GetSpeciesFromSlug(c)
	if species == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	if !IsWorldCreator(&species.World, user) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return species
}

func SocietyCreatorOnly(c *gin.Context) *models.Society {
	society := GetSocietyFromSlug(c)
	if society == nil {
		return nil
	}

	user := GetUserFromContext(c, true)
	if user == nil {
		return nil
	}

	if !IsWorldCreator(&society.World, user) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return nil
	}

	return society
}

func FilterSpeciesWorldAccess(species []models.Species, user *models.User) []models.Species {
	worldAccess := filterWorldAccess(species, user).([]models.Species)
	filtered := make([]models.Species, 0)
	for _, sp := range worldAccess {
		if sp.Public || IsWorldCreator(&sp.World, user) {
			filtered = append(filtered, sp)
		}
	}
	return filtered
}

func FilterSocietiesWorldAccess(societies []models.Society, user *models.User) []models.Society {
	worldAccess := filterWorldAccess(societies, user).([]models.Society)
	filtered := make([]models.Society, 0)
	for _, society := range worldAccess {
		if society.Public || IsWorldCreator(&society.World, user) {
			filtered = append(filtered, society)
		}
	}
	return filtered
}
