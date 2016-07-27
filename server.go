package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

var port = os.Getenv("PORT")

func init() {
	if port == "" {
		port = "3000"
	}
}

func server() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
	})

	router.PUT("/play-pause", func(c *gin.Context) {
	})

	router.PUT("/stop", func(c *gin.Context) {
	})

	router.PUT("/prev", func(c *gin.Context) {
	})

	router.PUT("/next", func(c *gin.Context) {
	})

	router.PUT("/repeat", func(c *gin.Context) {
	})

	router.PUT("/shuffle", func(c *gin.Context) {
	})

	router.Run(":" + port)
}
