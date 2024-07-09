package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/controllers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/middlewares"
	"log"
)

const apiVersion = "/v1"

func init() {
	initializers.ConnectToDB()
	initializers.MigrateDB()
}

func Root(c *gin.Context) {
	campaigns := map[string]string{
		"description": "Campaigns being played.",
		"collection":  "GET, POST*",
		"resource":    "GET, PUT*, DELETE*",
		"location":    apiVersion + "/campaigns",
	}

	characters := map[string]interface{}{
		"description": "Game characters (both PCs and NPCs).",
		"collection":  "GET, POST*",
		"resource":    "GET, PUT*, DELETE*",
		"location":    apiVersion + "/characters",
	}

	key := map[string]interface{}{
		"description": "Cycle your API key.",
		"collection":  "PUT*",
		"resource":    "",
		"location":    apiVersion + "/key",
	}

	rolls := map[string]interface{}{
		"description": "Rolls made on random tables.",
		"collection":  "GET, POST*",
		"resource":    "GET, DELETE*",
		"location":    apiVersion + "/rolls",
	}

	scrolls := map[string]interface{}{
		"description": "Scrolls and their current status.",
		"collection":  "GET, POST*",
		"resource":    "GET, PUT*, DELETE*",
		"location":    apiVersion + "/scrolls",
	}

	societies := map[string]interface{}{
		"description": "Fantasy societies and cultures.",
		"collection":  "GET, POST*",
		"resource":    "GET, PUT*, DELETE*",
		"location":    apiVersion + "/societies",
	}

	species := map[string]interface{}{
		"description": "Fantasy species (e.g., elves and goblins).",
		"collection":  "GET, POST*",
		"resource":    "GET, PUT*, DELETE*",
		"location":    apiVersion + "/species",
	}

	tables := map[string]interface{}{
		"description": "Random tables.",
		"collection":  "GET, POST*",
		"resource":    "GET, PUT*, DELETE*",
		"location":    apiVersion + "/tables",
	}

	users := map[string]interface{}{
		"description": "User records.",
		"collection":  "GET, POST*",
		"resource":    "GET, PUT*, DELETE*",
		"location":    apiVersion + "/users",
	}

	worlds := map[string]interface{}{
		"description": "Fictional worlds.",
		"collection":  "GET, POST*",
		"resource":    "GET, PUT*, DELETE*",
		"location":    apiVersion + "/worlds",
	}

	c.JSON(200, gin.H{
		"note":       "* These endpoints require authentication. Provide a valid API key for Bearer authentication.",
		"campaigns":  campaigns,
		"characters": characters,
		"key":        key,
		"rolls":      rolls,
		"scrolls":    scrolls,
		"societies":  societies,
		"species":    species,
		"tables":     tables,
		"users":      users,
		"worlds":     worlds,
	})
}

func main() {
	r := gin.Default()

	trustedProxies := []string{
		"127.0.0.1",
	}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	v := r.Group(apiVersion)
	{
		v.GET("/", Root)

		v.POST("/users", controllers.UserCreate)
		v.GET("/users", controllers.UserIndex)
		v.GET("/users/:username", controllers.UserRetrieve)

		authRequired := v.Group("/")
		authRequired.Use(middlewares.AuthRequired())
		{
			authRequired.PUT("/users", controllers.UserUpdate)
			authRequired.PUT("/key", controllers.KeyUpdate)

			authRequired.POST("/worlds", controllers.WorldCreate)
			authRequired.PUT("/worlds/:slug", controllers.WorldUpdate)
			authRequired.DELETE("/worlds/:slug", controllers.WorldDestroy)

			authRequired.POST("/campaigns", controllers.CampaignCreate)
			authRequired.PUT("/campaigns/:world/:slug", controllers.CampaignUpdate)
			authRequired.DELETE("/campaigns/:world/:slug", controllers.CampaignDestroy)

			authRequired.POST("/species", controllers.SpeciesCreate)
			authRequired.PUT("/species/:world/:slug", controllers.SpeciesUpdate)
			authRequired.DELETE("/species/:world/:slug", controllers.SpeciesDestroy)

			authRequired.POST("/societies", controllers.SocietyCreate)
			authRequired.PUT("/societies/:world/:slug", controllers.SocietyUpdate)
			authRequired.DELETE("/societies/:world/:slug", controllers.SocietyDestroy)

			authRequired.POST("/characters", controllers.CharacterCreate)
			authRequired.PUT("/characters/:id", controllers.CharacterUpdate)
			authRequired.DELETE("/characters/:id", controllers.CharacterDestroy)

			authRequired.POST("/scrolls", controllers.ScrollCreate)
			authRequired.PUT("/scrolls/:id", controllers.ScrollUpdate)
			authRequired.DELETE("/scrolls/:id", controllers.ScrollDestroy)
			authRequired.POST("/scrolls/:id/seal", controllers.ScrollSeal)
			authRequired.DELETE("/scrolls/:id/seal", controllers.ScrollUnseal)

			authRequired.POST("/tables", controllers.TableCreate)
			authRequired.PUT("/tables/:slug", controllers.TableUpdate)
			authRequired.DELETE("/tables/:slug", controllers.TableDestroy)

			authRequired.GET("/rolls", controllers.RollIndex)
			authRequired.GET("/rolls/:id", controllers.RollRetrieve)
			authRequired.DELETE("/rolls/:id", controllers.RollDestroy)
		}

		authOptional := v.Group("/")
		authOptional.Use(middlewares.AuthOptional())
		{
			authOptional.GET("/worlds", controllers.WorldIndex)
			authOptional.GET("/worlds/:slug", controllers.WorldRetrieve)

			authOptional.GET("/campaigns", controllers.CampaignIndex)
			authOptional.GET("/campaigns/:world/:slug", controllers.CampaignRetrieve)

			authOptional.GET("/species", controllers.SpeciesIndex)
			authOptional.GET("/species/:world/:slug", controllers.SpeciesRetrieve)

			authOptional.GET("/societies", controllers.SocietyIndex)
			authOptional.GET("/societies/:world/:slug", controllers.SocietyRetrieve)

			authOptional.GET("/characters", controllers.CharacterIndex)
			authOptional.GET("/characters/:id", controllers.CharacterRetrieve)

			authOptional.GET("/scrolls", controllers.ScrollIndex)
			authOptional.GET("/scrolls/:id", controllers.ScrollRetrieve)

			authOptional.GET("/tables", controllers.TableIndex)
			authOptional.GET("/tables/:slug", controllers.TableRetrieve)

			authOptional.POST("/rolls", controllers.RollCreate)
		}
	}

	err := r.Run()
	if err != nil {
		return
	}
}
