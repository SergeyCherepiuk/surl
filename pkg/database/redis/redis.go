package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var db *redis.Client

func MustConnect() {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	db := redis.NewClient(&redis.Options{Addr: addr, DB: 0})
	if err := db.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
}
