package parsers

import (
	"github.com/gin-gonic/gin"
	"time"
)

func BodyToKeyRequest(c *gin.Context) (string, string, string, bool) {
	var body struct {
		Username  string  `json:"username"`
		Password  string  `json:"password"`
		Ephemeral *bool   `json:"ephemeral"`
		Label     *string `json:"label"`
	}

	label := "Key created at " + time.Now().Format("2006-01-02 15:04:05")
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return "", "", label, true
	}

	ephemeral := true
	if body.Ephemeral != nil {
		ephemeral = *body.Ephemeral
	}

	if body.Label != nil {
		label = *body.Label
	}

	return body.Username, body.Password, label, ephemeral
}
