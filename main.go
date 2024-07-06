package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jefgodesky/rnrapi/controllers"
	"github.com/jefgodesky/rnrapi/initializers"
)

const apiVersion = "/v1"

func init() {
	initializers.ConnectToDB()
	initializers.MigrateDB()
}

func main() {
	r := gin.Default()
	v1 := r.Group(apiVersion)
	{
		v1.POST("/users", controllers.UserCreate)
		v1.GET("/users", controllers.UserIndex)
		v1.GET("/users/:username", controllers.UserRetrieve)
		v1.PUT("/users/:username", controllers.UserUpdate)
	}
	r.Run()
}
