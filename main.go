package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stewart/cmus"
)

var port = os.Getenv("PORT")
var router = gin.Default()
var client = cmus.Client{}

func init() {
	if port == "" {
		port = "3000"
	}

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	router.GET("/", func(c *gin.Context) {
		status, err := client.Status()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(200, gin.H{
			"playing":  status.Playing,
			"file":     status.File,
			"duration": status.Duration,
			"position": status.Position,
			"tags":     status.Tags,
			"settings": status.Settings,
		})
	})

	router.Run(":" + port)
}
