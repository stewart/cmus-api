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

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

	router.Run(":" + port)
}
