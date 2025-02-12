package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// ConnectRedis establishes a connection to the Redis server.
func ConnectRedis() error {
	addr := os.Getenv("REDIS_ADDR") // Get address from environment variable
	if addr == "" {
		return fmt.Errorf("REDIS_ADDR environment variable not set")
	}

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           0,                          
		DialTimeout:  5 * time.Second,            
		ReadTimeout:  30 * time.Second,           
		WriteTimeout: 30 * time.Second,          
	})


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
    if err != nil {
        return fmt.Errorf("failed to ping redis: %w", err)
    }

	redisClient = client
	log.Println("Connected to Redis successfully")
	return nil
}

// GetRedis returns the Redis client instance.  This will return nil if the connection failed.
func GetRedis() *redis.Client {
	return redisClient
}

