package cache

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// New creates a new Redis client and it can talk to the server
func New(address string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		slog.With("package", "cache").Error(err.Error())
		return nil, err
	}

	return rdb, nil
}
