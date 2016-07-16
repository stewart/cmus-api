package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var port = os.Getenv("PORT")

func init() {
	if port == "" {
		port = "3000"
	}
}

func main() {
	router := gin.Default()

	if err := state.cmus.Connect(); err != nil {
		log.Fatal(err)
	}

	status, err := state.cmus.Status()
	if err != nil {
		log.Fatal(err)
	}

	state.Lock()
	state.status = status
	state.err = nil
	state.Unlock()

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
		command(c, state.cmus.PlayPause())
	})

	router.PUT("/stop", func(c *gin.Context) {
		command(c, state.cmus.Stop())
	})

	router.PUT("/prev", func(c *gin.Context) {
		command(c, state.cmus.Prev())
	})

	router.PUT("/next", func(c *gin.Context) {
		command(c, state.cmus.Next())
	})

	router.PUT("/repeat", func(c *gin.Context) {
		command(c, state.cmus.Repeat())
	})

	router.PUT("/shuffle", func(c *gin.Context) {
		command(c, state.cmus.Shuffle())
	})

	loop()

	router.Run(":" + port)
}
