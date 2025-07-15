package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"thanhldt060802/redisclient"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var EXAMPLE_NUM = 2
var EXAMPLES map[int]func()

func init() {
	EXAMPLES = map[int]func(){
		1: Example1,
		2: Example2,
	}
}

func main() {

	EXAMPLES[EXAMPLE_NUM]()

}

/*
- Example for Subscribe() and Publish()

- Subscribe() will listen for Payloads sent over TCP Socket

- Publish() will send Payloads to TCP Socket
*/
func Example1() {
	redisClient := redisclient.NewRedisClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "12345678",
		DB:       0,
	})
	redisPub := redisclient.NewRedisPub[string](redisClient)
	redisSub := redisclient.NewRedisSub[string](redisClient)

	go func() {
		count := 0

		redisSub.Subscribe(context.Background(), "my-channel", func(data string) {
			fmt.Println(data)

			count++
			fmt.Printf("Count: %v\n", count)

			time.Sleep(1 * time.Second)
		})
	}()

	go func() {
		for i := 1; ; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			data := fmt.Sprintf("my-payload-%v", i)
			if err := redisPub.Publish(ctx, "my-channel", data); err != nil {
				cancel()
				return
			}

			cancel()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {}
}

/*
- Ref: Example2()

- Using data struct
*/
func Example2() {
	type SubDataStruct struct {
		Field1 string `json:"field1"`
		Field2 int32  `json:"field2"`
		Field3 int64  `json:"field3"`
	}
	type DataStruct struct {
		Field1 string        `json:"field1"`
		Field2 int32         `json:"field2"`
		Field3 int64         `json:"field3"`
		Field4 float32       `json:"field4"`
		Field5 float64       `json:"field5"`
		Field6 time.Time     `json:"field6"`
		Field7 SubDataStruct `json:"field7"`
	}

	redisClient := redisclient.NewRedisClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "12345678",
		DB:       0,
	})
	redisPub := redisclient.NewRedisPub[*DataStruct](redisClient)
	redisSub := redisclient.NewRedisSub[*DataStruct](redisClient)

	go func() {
		count := 0

		redisSub.Subscribe(context.Background(), "my-channel", func(data *DataStruct) {
			fmt.Println(*data)

			count++
			fmt.Printf("Count: %v\n", count)

			time.Sleep(1 * time.Second)
		})
	}()

	go func() {
		for i := 1; ; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			data := DataStruct{
				Field1: uuid.New().String(),
				Field2: rand.Int32(),
				Field3: rand.Int64(),
				Field4: rand.Float32(),
				Field5: rand.Float64(),
				Field6: time.Now(),
				Field7: SubDataStruct{
					Field1: uuid.New().String(),
					Field2: rand.Int32(),
					Field3: rand.Int64(),
				},
			}
			if err := redisPub.Publish(ctx, "my-channel", &data); err != nil {
				cancel()
				return
			}

			cancel()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {}
}
