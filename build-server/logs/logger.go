package logs

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASS")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	if redisPass == "" {
		redisPass = ""
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		log.Printf("Connected to Redis at %s", redisAddr)
	}
}

func PublishLog(channel, message string) {
	err := redisClient.Publish(context.Background(), channel, message).Err()
	if err != nil {
		log.Printf("Failed to publish log to Redis: %v", err)
	}
}
