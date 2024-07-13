package middlewares

import "github.com/gin-gonic/gin"

var AllowedMethods = map[string]string{
	"/v1/":                        "GET, HEAD",
	"/v1/campaigns":               "GET, HEAD, POST",
	"/v1/campaigns/:world/:slug":  "GET, HEAD, PUT, DELETE",
	"/v1/characters":              "GET, HEAD, POST",
	"/v1/characters/:id":          "GET, HEAD, PUT, DELETE",
	"/v1/emails":                  "GET, HEAD, POST",
	"/v1/emails/:id":              "GET, HEAD, DELETE",
	"/v1/emails/:id/verification": "POST",
	"/v1/keys":                    "GET, HEAD, POST",
	"/v1/keys/:id":                "GET, HEAD, DELETE",
	"/v1/rolls":                   "GET, HEAD",
	"/v1/rolls/:id":               "GET, HEAD, DELETE",
	"/v1/scrolls":                 "GET, HEAD",
	"/v1/scrolls/:id":             "GET, HEAD, PUT, DELETE",
	"/v1/societies":               "GET, HEAD",
	"/v1/societies/:world/:slug":  "GET, HEAD, PUT, DELETE",
	"/v1/species":                 "GET, HEAD",
	"/v1/species/:world/:slug":    "GET, HEAD, PUT, DELETE",
	"/v1/tables":                  "GET, HEAD",
	"/v1/tables/:slug":            "GET, HEAD, PUT, DELETE",
	"/v1/users":                   "GET, HEAD, POST, PUT, DELETE",
	"/v1/users/:username":         "GET, HEAD",
	"/v1/worlds":                  "GET, HEAD",
	"/v1/worlds/:slug":            "GET, HEAD, PUT, DELETE",
}

func AllowHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		methods, exists := AllowedMethods[path]
		if exists {
			c.Header("Allow", methods)
			c.Header("Access-Control-Allow-Methods", methods)
		}
		c.Next()
	}
}
