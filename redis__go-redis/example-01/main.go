package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"thanhldt060802/redisclient"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var EXAMPLE_NUM = 2
var EXAMPLES map[int]func()

var RedisClient redisclient.IRedisClient

func init() {
	EXAMPLES = map[int]func(){
		1: Example1,
		2: Example2,
	}

	RedisClient = redisclient.NewRedisClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "12345678",
		DB:       0,
	})
}

func main() {

	EXAMPLES[EXAMPLE_NUM]()

}

/*
- Example for Subscribe() and Publish()
- Subscribe() will listen for Payloads sent over TCP Socket
- Publish() will send Payloads to TCP Socket
- It is so simple if we use with string
*/
func Example1() {
	go func() {
		count := 0

		RedisClient.Subscribe(context.Background(), "my-channel", func(payload string) {
			data := payload
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
			if err := RedisClient.Publish(ctx, "my-channel", data); err != nil {
				log.Infof("Publish %v to my-channel failed: %v", data, err.Error())
				continue
			}
			log.Infof("Publish %v to my-channel successful", data)

			cancel()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {}
}

/*
- Example for Subscribe() and Publish()
- When using data struct, it should be passed to Marshal() before Publish() and from Unmarshal() to receive data after Subscribe()
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

	go func() {
		count := 0

		RedisClient.Subscribe(context.Background(), "my-channel", func(payload string) {
			var data DataStruct
			if err := json.Unmarshal([]byte(payload), &data); err != nil {
				log.Infof("Unmarshal %v failed: %v", payload, err.Error())
				return
			}
			fmt.Println(data)

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
			payload, err := json.Marshal(data)
			if err != nil {
				log.Infof("Marshal data failed: %v", err.Error())
				continue
			}
			if err := RedisClient.Publish(ctx, "my-channel", payload); err != nil {
				log.Infof("Publish %v to my-channel failed: %v", data, err.Error())
				continue
			}
			log.Infof("Publish %v to my-channel successful", data)

			cancel()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {}
}
