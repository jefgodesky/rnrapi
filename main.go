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
	c.JSON(200, middlewares.AllowedMethods)
}

func main() {
	r := gin.Default()
	r.Use(middlewares.AllowHeaderMiddleware())

	trustedProxies := []string{
		"127.0.0.1",
	}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	v := r.Group(apiVersion)
	{
		v.GET("/", Root)
		v.HEAD("/", Root)

		v.POST("/users", controllers.UserCreate)
		v.GET("/users", controllers.UserIndex)
		v.HEAD("/users", controllers.UserIndex)
		v.GET("/users/:username", controllers.UserRetrieve)
		v.HEAD("/users/:username", controllers.UserRetrieve)
		v.POST("/keys", controllers.KeyCreate)
		v.POST("/password-reset", controllers.PasswordReset)

		authRequired := v.Group("/")
		authRequired.Use(middlewares.AuthRequired())
		{
			authRequired.PUT("/users", controllers.UserUpdate)
			authRequired.DELETE("/users", controllers.UserDestroy)

			authRequired.GET("/keys", controllers.KeyIndex)
			authRequired.HEAD("/keys", controllers.KeyIndex)
			authRequired.GET("/keys/:id", controllers.KeyRetrieve)
			authRequired.HEAD("/keys/:id", controllers.KeyRetrieve)
			authRequired.DELETE("/keys/:id", controllers.KeyDestroy)

			authRequired.POST("/emails", controllers.EmailCreate)
			authRequired.HEAD("/emails", controllers.EmailIndex)
			authRequired.GET("/emails", controllers.EmailIndex)
			authRequired.HEAD("/emails/:id", controllers.EmailRetrieve)
			authRequired.GET("/emails/:id", controllers.EmailRetrieve)
			authRequired.DELETE("/emails/:id", controllers.EmailDestroy)
			authRequired.POST("/emails/:id/verification", controllers.EmailVerify)

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
			authRequired.HEAD("/rolls", controllers.RollIndex)
			authRequired.GET("/rolls/:id", controllers.RollRetrieve)
			authRequired.HEAD("/rolls/:id", controllers.RollRetrieve)
			authRequired.DELETE("/rolls/:id", controllers.RollDestroy)
		}

		authOptional := v.Group("/")
		authOptional.Use(middlewares.AuthOptional())
		{
			authOptional.GET("/worlds", controllers.WorldIndex)
			authOptional.HEAD("/worlds", controllers.WorldIndex)
			authOptional.GET("/worlds/:slug", controllers.WorldRetrieve)
			authOptional.HEAD("/worlds/:slug", controllers.WorldRetrieve)

			authOptional.GET("/campaigns", controllers.CampaignIndex)
			authOptional.HEAD("/campaigns", controllers.CampaignIndex)
			authOptional.GET("/campaigns/:world/:slug", controllers.CampaignRetrieve)
			authOptional.HEAD("/campaigns/:world/:slug", controllers.CampaignRetrieve)

			authOptional.GET("/species", controllers.SpeciesIndex)
			authOptional.HEAD("/species", controllers.SpeciesIndex)
			authOptional.GET("/species/:world/:slug", controllers.SpeciesRetrieve)
			authOptional.HEAD("/species/:world/:slug", controllers.SpeciesRetrieve)

			authOptional.GET("/societies", controllers.SocietyIndex)
			authOptional.HEAD("/societies", controllers.SocietyIndex)
			authOptional.GET("/societies/:world/:slug", controllers.SocietyRetrieve)
			authOptional.HEAD("/societies/:world/:slug", controllers.SocietyRetrieve)

			authOptional.GET("/characters", controllers.CharacterIndex)
			authOptional.HEAD("/characters", controllers.CharacterIndex)
			authOptional.GET("/characters/:id", controllers.CharacterRetrieve)
			authOptional.HEAD("/characters/:id", controllers.CharacterRetrieve)

			authOptional.GET("/scrolls", controllers.ScrollIndex)
			authOptional.HEAD("/scrolls", controllers.ScrollIndex)
			authOptional.GET("/scrolls/:id", controllers.ScrollRetrieve)
			authOptional.HEAD("/scrolls/:id", controllers.ScrollRetrieve)

			authOptional.GET("/tables", controllers.TableIndex)
			authOptional.HEAD("/tables", controllers.TableIndex)
			authOptional.GET("/tables/:slug", controllers.TableRetrieve)
			authOptional.HEAD("/tables/:slug", controllers.TableRetrieve)

			authOptional.POST("/rolls", controllers.RollCreate)
		}
	}

	err := r.Run()
	if err != nil {
		return
	}
}
