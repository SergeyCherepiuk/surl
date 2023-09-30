package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	sessionsDb *redis.Client
	cacheDb    *redis.Client
)

func MustConnect() {
	sessionsAddr := fmt.Sprintf("%s:%s", os.Getenv("SESSIONS_REDIS_HOST"), os.Getenv("SESSIONS_REDIS_PORT"))
	fmt.Println(sessionsAddr)
	sessionsDb = redis.NewClient(&redis.Options{Addr: sessionsAddr, DB: 0})
	if err := sessionsDb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	cacheAddr := fmt.Sprintf("%s:%s", os.Getenv("CACHE_REDIS_HOST"), os.Getenv("CACHE_REDIS_PORT"))
	fmt.Println(cacheAddr)
	cacheDb = redis.NewClient(&redis.Options{Addr: cacheAddr, DB: 0})
	if err := cacheDb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
}
