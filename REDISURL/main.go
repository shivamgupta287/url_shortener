package main

import (
	"fmt"
	"redisurl/handler"
	"redisurl/store"

	"github.com/gin-gonic/gin"
)

func main() {
  fmt.Println("Hello world!")
  r := gin.Default()
  r.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "Hey Go URL shortener",
    })
  })

  r.POST("/create-short-url", func(c *gin.Context) {
    handler.CreateShortUrl(c)
  })

  r.GET("/:shortUrl", func(c *gin.Context) {
    handler.HandleShortUrlRedirect(c)
  })

  store.InitializeStore()

  err := r.Run(":5000")
  if err != nil {
    panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
  }
}