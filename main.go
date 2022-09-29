package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var ctx = context.Background()

func GetCountPattern(rdb *redis.Client, pattern string) {
	counter := 0
	iter := rdb.Scan(ctx, 0, pattern+":*", 0).Iterator()
	for iter.Next(ctx) {

		value, err := rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			panic(err)
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
		panic(err)
	}
	fmt.Println(counter)
}

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	t, err := rdb.Incr(ctx, "k:1").Result()
	if err != nil {
		panic(err)
	}

	_, err = rdb.Expire(ctx, "k:1", 4*time.Second).Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v, %T\n", t, t)

	GetCountPattern(rdb, "k")
	GetCountPattern(rdb, "c")
}

func main() {
	ExampleClient()
}
