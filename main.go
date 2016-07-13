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

	command := func(c *gin.Context, err error) {
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(204)
	}

	router.PUT("/play-pause", func(c *gin.Context) {
		command(c, client.PlayPause())
	})

	router.PUT("/stop", func(c *gin.Context) {
		command(c, client.Stop())
	})

	router.PUT("/prev", func(c *gin.Context) {
		command(c, client.Prev())
	})

	router.PUT("/next", func(c *gin.Context) {
		command(c, client.Next())
	})

	router.PUT("/repeat", func(c *gin.Context) {
		command(c, client.Repeat())
	})

	router.PUT("/shuffle", func(c *gin.Context) {
		command(c, client.Shuffle())
	})

	loop()

	router.Run(":" + port)
}
