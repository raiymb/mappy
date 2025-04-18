package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, addr, pwd string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    pwd,
		DB:          db,
		DialTimeout: 5 * time.Second,
	})
	return rdb, rdb.Ping(ctx).Err()
}
