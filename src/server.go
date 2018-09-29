package main

import "fmt"
import "github.com/go-redis/redis"
import (
	"github.com/jinzhu/gorm"
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

func main() {
	redisClient := RedisClient()

	pong, err := redisClient.Ping().Result()
	fmt.Println("Redis ping", pong, err)

	dbClient := PostgresClient()
	defer dbClient.Close()

	dbClient.AutoMigrate(&urls{})

	dbClient.Create(&urls{Tinyurl: "test.tiny", Longurl: "test.t"})

	dbClient.Close()
}
