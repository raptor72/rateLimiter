package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/raptor72/rateLimiter/config"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var ctx = context.Background()

type Client struct {
	cfg *config.Config
	ctx context.Context
	rdb *redis.Client
}

func NewClient(config *config.Config) *Client {
	//ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	return &Client{
		cfg: config,
		ctx: ctx,
		rdb: rdb,
	}
}

func (c *Client) GetCountPattern(pattern string) (*int, error) {
	counter := 0

	iter := c.rdb.Scan(c.ctx, 0, pattern+":*", 0).Iterator()
	for iter.Next(c.ctx) {

		value, err := c.rdb.Get(c.ctx, iter.Val()).Result()
		if err != nil {
			return nil, err
		}

		var numValue int

		numValue, err = strconv.Atoi(value)
		if err != nil {
            log.Warn("Wrong incremented value")
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

func (c *Client) IncrementOrBlock(key string, limitCount int, ttl time.Duration) bool {
	now := time.Now()
	sec := now.Unix()

	strSec := strconv.Itoa(int(sec))

	fullKey := key + ":" + strSec

	count, err := c.GetCountPattern(fullKey)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if *count <= limitCount {
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
