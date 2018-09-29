package main

import (
	"context"
	"fmt"
	"io"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"net/http"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// urls -> URL database structure
type urls struct {
	gorm.Model
	Tinyurl string
	Longurl string
}

// PostgresClient -> Provides a connection to the postgres database server
func PostgresClient() *gorm.DB {
	dbClient, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=tiny_scale_go password=<db password> sslmode=disable")
	if err != nil {
		panic(err)
	}
	return dbClient
}

// RedisClient -> Provides a connection to the Redis server
func RedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

// IndexHandler -> Handles requests coming to / route
func IndexHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Welcome!\n")
}

// StopHandler -> Stops the server on request to /stop route
func StopHandler(w http.ResponseWriter, req *http.Request, dbClient *gorm.DB, redisClient *redis.Client, serverInstance *http.Server) {
	fmt.Println("Stopping server...\n")
	dbClient.Close()
	redisClient.Close()
	serverInstance.Shutdown(context.TODO())
}

func main() {
	redisClient := RedisClient()

	pong, err := redisClient.Ping().Result()
	fmt.Println("Redis ping", pong, err)

	dbClient := PostgresClient()
	defer dbClient.Close()

	dbClient.AutoMigrate(&urls{})

	serverInstance := &http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", IndexHandler)

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		StopHandler(w, r, dbClient, redisClient, serverInstance)
	})

	serverInstance.ListenAndServe()
}
