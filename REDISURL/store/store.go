package store

import (
	"context"
	"fmt"
	"redisurl/db"

	"github.com/go-redis/redis/v8"
)

type StorageService struct {
  redisClient *redis.Client
}

var (
  storeService = &StorageService{}
  ctx = context.Background()
)

func InitializeStore() *StorageService {
  redisClient := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
  })

  pong, err := redisClient.Ping(ctx).Result()
  if err != nil {
    panic(fmt.Sprintf("Error init Redis: %v", err))
  }

  fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)
  storeService.redisClient = redisClient
  return storeService
}

func SaveUrlMapping(shortUrl string, originalUrl string) {
  db.SaveTOMongo(shortUrl,originalUrl)
  err := storeService.redisClient.Set(ctx,originalUrl,shortUrl,0).Err()
  if err != nil {
    panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
  }
}

func RetrieveInitialUrl(shortUrl string) string {
  result, err := storeService.redisClient.Get(ctx, shortUrl).Result()
  if err != nil {
    return db.RetrieveInitialUrlMongo(shortUrl)
  }
  return result
}

func CheckFromRedis(longurl string) (bool,string) {
  exists, err := storeService.redisClient.Exists(ctx, longurl).Result()
  if err != nil {
    panic(err)
  }
  if exists == 1 {
    result , _ := storeService.redisClient.Get(ctx,longurl).Result()
    db.SaveTOMongo(result,longurl)
    return true,result
    } else {
      return false,"Key does not exist"
    }
}