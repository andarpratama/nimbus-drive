package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
}
