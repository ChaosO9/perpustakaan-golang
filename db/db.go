package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	// _redis "github.com/go-redis/redis/v7"
	_ "github.com/lib/pq" //import postgres
)


var db *sqlx.DB

func Connect() (*sqlx.DB, error) {
	host := os.Getenv("POSTGRES_HOST") // Get the host from environment variables
	if host == "" {
		host = "localhost" //Default to localhost if not specified.
	}
	port := os.Getenv("POSTGRES_PORT") // Get the port from environment variables
	if port == "" {
		port = "5432" //Default to 5432 if not specified
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

func SetDB(database *sqlx.DB) {
    db = database
}

func GetDB() *sqlx.DB {
    return db
}

//RedisClient ...
// var RedisClient *_redis.Client

// //InitRedis ...
// func InitRedis(selectDB ...int) {

// 	var redisHost = os.Getenv("REDIS_HOST")
// 	var redisPassword = os.Getenv("REDIS_PASSWORD")

// 	RedisClient = _redis.NewClient(&_redis.Options{
// 		Addr:     redisHost,
// 		Password: redisPassword,
// 		DB:       selectDB[0],
// 		// DialTimeout:        10 * time.Second,
// 		// ReadTimeout:        30 * time.Second,
// 		// WriteTimeout:       30 * time.Second,
// 		// PoolSize:           10,
// 		// PoolTimeout:        30 * time.Second,
// 		// IdleTimeout:        500 * time.Millisecond,
// 		// IdleCheckFrequency: 500 * time.Millisecond,
// 		// TLSConfig: &tls.Config{
// 		// 	InsecureSkipVerify: true,
// 		// },
// 	})

// }

// //GetRedis ...
// func GetRedis() *_redis.Client {
// 	return RedisClient
// }