package db

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveTOMongo(shortUrl string, originalUrl string) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer client.Disconnect(context.Background())
	collection := client.Database("REDIS_URL_SHORTENER").Collection("urls")
	filter := bson.M{"originalurl":originalUrl}
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		panic(err)
	}
	if count < 1 {
		_, insertErr := collection.InsertOne(ctx, bson.D{{Key: "originalurl", Value: originalUrl}, {Key: "shorturl", Value: shortUrl}})
		if insertErr != nil {
			log.Fatal(insertErr)
		}
	}
}


func RetrieveInitialUrlMongo(shortUrl string) string {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel:= context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer client.Disconnect(context.Background())
	collection := client.Database("REDIS_URL_SHORTENER").Collection("urls")
	var result bson.M
	err = collection.FindOne(ctx, bson.M{"shorturl": shortUrl}).Decode(&result)
	if err != nil {
		log.Println(err)
	}
	SaveToRedis(shortUrl,result["originalurl"].(string))
	return result["originalurl"].(string)
}


func SaveToRedis(shortUrl string, originalUrl string)  {
	ctx, cancel:= context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
    	Password: "",
   		DB:       0,
	})
	ERR:= redisClient.Set(ctx, shortUrl,originalUrl,0).Err()
	if ERR != nil {
		panic(ERR)
	}
}



func CheckFromMongo(longurl string) (bool,string) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel:= context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer client.Disconnect(context.Background())
	collection := client.Database("REDIS_URL_SHORTENER").Collection("urls")
	var result bson.M
	err = collection.FindOne(ctx, bson.M{"originalurl": longurl}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false,"Longurl not found in db"
		} else if err != nil {
			panic(err)
		}
		ThenSaveToRedis(longurl,result["shorturl"].(string))
		return true,result["shorturl"].(string)
}


func ThenSaveToRedis(originalUrl string ,shortUrl string)  {
	ctx, cancel:= context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
    	Password: "",
   		DB:       0,
	})
	ERR:= redisClient.Set(ctx,originalUrl,shortUrl,0).Err()
	if ERR != nil {
		panic(ERR)
	}
}