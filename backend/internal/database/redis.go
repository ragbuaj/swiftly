package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis initializes the Redis connection
func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // No password by default in dev
		DB:       0,  // Use default DB
	})

	// Test connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis")
}

// GetRedis returns the Redis client instance
func GetRedis() *redis.Client {
	return redisClient
}

// CloseRedis closes the Redis connection
func CloseRedis() {
	if redisClient != nil {
		redisClient.Close()
	}
}
