package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/controllers"
	"github.com/jefgodesky/rnrapi/initializers"
	"github.com/jefgodesky/rnrapi/middlewares"
)

const apiVersion = "/v1"

func init() {
	initializers.ConnectToDB()
	initializers.MigrateDB()
}

func main() {
	r := gin.Default()
	v := r.Group(apiVersion)
	{
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
		}
	}

	err := r.Run()
	if err != nil {
		return
	}
}
