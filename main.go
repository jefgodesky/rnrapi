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
		}

		authOptional := v.Group("/")
		authOptional.Use(middlewares.AuthOptional())
		{
			authOptional.GET("/worlds", controllers.WorldIndex)
			authOptional.GET("/worlds/:slug", controllers.WorldRetrieve)
		}
	}
	r.Run()
}
