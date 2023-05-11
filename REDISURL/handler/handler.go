package handler

import (
	"fmt"
	"net/http"
	"redisurl/db"
	"redisurl/shortener"
	"redisurl/store"

	"github.com/gin-gonic/gin"
)

// Request model definition
type UrlCreationRequest struct {
  LongUrl string `json:"long_url" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
  var creationRequest UrlCreationRequest
  if err := c.ShouldBindJSON(&creationRequest); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }


  flag , url := store.CheckFromRedis("https://"+creationRequest.LongUrl)
  if flag {
    host := fmt.Sprintf("http://%s/", c.Request.Host)
    c.JSON(200, gin.H{
    "message": "short url from redis",
    "short_url": host+url,
  })
  } else {
    flag2,url2 := db.CheckFromMongo("https://"+creationRequest.LongUrl)
    if flag2 {
      host := fmt.Sprintf("http://%s/", c.Request.Host)
      c.JSON(200, gin.H{
        "message": "short url from mongo",
        "short_url": host+url2,
      })
    } else {
      shortUrl := shortener.GenerateShortLink(creationRequest.LongUrl)
      store.SaveUrlMapping(shortUrl, "https://"+creationRequest.LongUrl)
      host := fmt.Sprintf("http://%s/", c.Request.Host)
      c.JSON(200, gin.H{
        "message": "short url created successfully",
        "short_url": host + shortUrl,
      })
    }
  }
}

func HandleShortUrlRedirect(c *gin.Context) {
  shortUrl := c.Param("shortUrl")
  initialUrl := store.RetrieveInitialUrl(shortUrl)
  c.Redirect(302, initialUrl)
}