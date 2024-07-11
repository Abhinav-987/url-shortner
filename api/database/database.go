package database

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	Ctx    = context.Background()
	Client *redis.Client
)

func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})
	log.Println("Redis is Connected")
	return rdb
}

func InitializeClient() {
	Client = CreateClient(0)
}
