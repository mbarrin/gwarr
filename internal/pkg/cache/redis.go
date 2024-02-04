package cache

import (
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		slog.With("package", "cache").Error(err.Error())
		os.Exit(1)
	}

	return rdb
}
