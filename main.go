package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stewart/cmus"
)

var port = os.Getenv("PORT")
var client = cmus.Client{}

var state = &State{}

func init() {
	if port == "" {
		port = "3000"
	}

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	status, err := client.Status()
	if err != nil {
		log.Fatal(err)
	}

	state.Lock()
	state.status = status
	state.err = nil
	state.Unlock()
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		state.RLock()
		defer state.RUnlock()

		if state.err != nil {
			c.JSON(500, gin.H{"error": state.err.Error()})
			return
		}

		status := state.status

		c.JSON(200, gin.H{
			"playing":  status.Playing,
			"file":     status.File,
			"duration": status.Duration,
			"position": status.Position,
			"tags":     status.Tags,
			"settings": status.Settings,
		})
	})

	loop()

	router.Run(":" + port)
}
