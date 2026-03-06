package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ConnectCache(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // может быть redis:6379
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, err
}
