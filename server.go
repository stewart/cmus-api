package main

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stewart/cmus"
)

var port = os.Getenv("PORT")

func init() {
	if port == "" {
		port = "3000"
	}
}

func server() {
	router := gin.Default()

	router.Use(CORS())

	command := func(c *gin.Context, err error) {
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(204)
	}

	router.GET("/", func(c *gin.Context) {
		state.RLock()
		defer state.RUnlock()

		if state.err != nil {
			c.JSON(500, gin.H{
				"error": state.err.Error(),
			})
			return
		}

		c.JSON(200, serializeStatus(state.status))
	})

	router.GET("/sse", func(c *gin.Context) {
		initial := false
		prevStatus := &cmus.Status{}
		var prevErr error

		c.Stream(func(w io.Writer) bool {
			state.RLock()
			defer state.RUnlock()

			status := state.status
			err := state.err

			if err == nil {
				// if previous message was an error, or status has changed, send
				if prevErr != nil || isDifferentStatus(prevStatus, status) || !initial {
					c.SSEvent("status", serializeStatus(status))
					initial = true
				}
			} else {
				// if a new error, send error
				if prevErr == nil || prevErr.Error() != err.Error() {
					c.SSEvent("error", err.Error())
				}
			}

			prevStatus = status
			prevErr = err

			time.Sleep(50 * time.Millisecond)

			return true
		})
	})

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

	router.PUT("/volume", func(c *gin.Context) {
		post := struct{ Level string }{}

		c.BindJSON(&post)

		if post.Level == "" {
			c.JSON(400, gin.H{"error": "volume level not supplied"})
			return
		}

		command(c, client.Volume(post.Level))
	})

	router.PUT("/seek", func(c *gin.Context) {
		post := struct{ Position string }{}

		c.BindJSON(&post)

		if post.Position == "" {
			c.JSON(400, gin.H{"error": "seek position not supplied"})
			return
		}

		command(c, client.Seek(post.Position))
	})

	router.Run(":" + port)
}
