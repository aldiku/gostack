package config

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

// RC Context
var RC *redis.Client

// Redis Initialization
func Redis() {
	addr := os.Getenv("CACHE_URL")
	pass := os.Getenv("CACHE_PASSWORD")

	RC = redis.NewClient(&redis.Options{Addr: addr, Password: pass, DB: 0})

	_, err := RC.Ping().Result()

	if err != nil {
		log.Panic(err)
	}
	log.Println("Redis Connected")
}
