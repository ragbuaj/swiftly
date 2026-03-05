package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	var opts *redis.Options
	var err error

	// Check if it's a URL format (redis:// or rediss://)
	if strings.HasPrefix(redisURL, "redis://") || strings.HasPrefix(redisURL, "rediss://") {
		opts, err = redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("Invalid Redis URL: %v", err)
		}
	} else {
		// Fallback to old behavior for raw host:port
		opts = &redis.Options{
			Addr: redisURL,
		}
	}

	redisClient = redis.NewClient(opts)

	// Test connection
	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis")
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	if redisClient != nil {
		redisClient.Close()
	}
}
