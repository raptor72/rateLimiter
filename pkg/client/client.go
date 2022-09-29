package client

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/raptor72/rateLimiter/config"
)

type Client struct {
	cfg *config.Config
	ctx context.Context
	rdb *redis.Client
}

func NewClient(config *config.Config, context context.Context) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	return &Client{
		cfg: config,
		ctx: context,
		rdb: rdb,
	}
}
