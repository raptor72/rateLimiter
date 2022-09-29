package limiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/raptor72/rateLimiter/config"
	"github.com/raptor72/rateLimiter/pkg/client"
	"log"
	"strconv"
	"time"
)

var ctx = context.Background()

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

func NewClients() *client.Client {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	cClient := client.NewClient(cfg, ctx)
	return cClient
}

func (c *Client) GetCountPattern(pattern string) (*int, error) {
	counter := 0

	iter := c.rdb.Scan(c.ctx, 0, pattern+":*", 0).Iterator()
	for iter.Next(ctx) {

		value, err := c.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			return nil, err
		}

		var numValue int

		numValue, err = strconv.Atoi(value)
		if err != nil {
			fmt.Println("Probably wrong")
			continue
		}
		if numValue != 0 && numValue != 1 {
			counter += numValue
		} else {
			counter++
		}

	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return &counter, nil
}

func (c *Client) IncrementOrBlock(key string, limit int, ttl time.Duration) bool {
	now := time.Now()
	sec := now.Unix()

	strSec := strconv.Itoa(int(sec))

	fullKey := key + ":" + strSec

	count, err := c.GetCountPattern(fullKey)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if *count <= limit {
		err = c.IncrementWithTTL(fullKey, ttl)
		if err != nil {
			log.Fatal(err)
			return false
		}
		return true
	}

	return false
}

func (c *Client) IncrementWithTTL(key string, sec time.Duration) error {
	_, err := c.rdb.Incr(c.ctx, key).Result()
	if err != nil {
		return err
	}

	_, err = c.rdb.Expire(c.ctx, key, sec*time.Second).Result()
	if err != nil {
		return err
	}
	return nil
}