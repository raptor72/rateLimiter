package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func GetCountPattern(rdb *redis.Client, pattern string) {
    counter := 0
    iter := rdb.Scan(ctx, 0,  pattern + ":*", 0).Iterator()
	for iter.Next(ctx) {
		// fmt.Println("keys", iter.Val())
        counter++
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

	// fmt.Printf("rdb %v, type %T\n", rdb, rdb) // Redis<localhost:6379 db:0>, type *redis.Client

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist

	set, err := rdb.SetNX(ctx, "k:123", "value", 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(set)

	set, err = rdb.SetNX(ctx, "k:3egwwe", "value", 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(set)


    t, err := rdb.Incr(ctx, "k:3egweu").Result()
    if err != nil {
		panic(err)
	}

    r, err := rdb.Expire(ctx, "k:3egweu", 3*time.Second).Result()
    // rdb.Expire()
    if err != nil {
		panic(err)
	}
    fmt.Printf("%v\n", r)

    fmt.Printf("%v, %T\n",t, t)


    GetCountPattern(rdb, "k")
    GetCountPattern(rdb, "c")
}

func main() {
	ExampleClient()
}
