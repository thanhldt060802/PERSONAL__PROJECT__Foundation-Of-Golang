package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"thanhldt060802/common/pubsub"
	"thanhldt060802/internal/rabbitmqclient"
	"thanhldt060802/model"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

func Example1() {
	rabbitmqclient.RabbitMQClientConnInstance = rabbitmqclient.NewRabbitMQClient(rabbitmqclient.RabbitMQConfig{
		Host:     "pitel-cx.dev.tel4vn.com",
		Port:     5672,
		Username: "pitel",
		Password: "Pitel@8229",
	})
	if err := rabbitmqclient.RabbitMQClientConnInstance.DeclareExchange("user", "topic"); err != nil {
		log.Fatalf("Failed to declare new exchange: %v", err.Error())
	}
	if rabbitMqPub, err := pubsub.NewRabbitMQPub[string](); err != nil {
		log.Fatalf("Failed to create new publisher: %v", err.Error())
	} else {
		pubsub.RabbitMQPubInstance1 = rabbitMqPub
	}
	if rabbitMqSub, err := pubsub.NewRabbitMQSub[string](); err != nil {
		log.Fatalf("Failed to create new subscriber: %v", err.Error())
	} else {
		pubsub.RabbitMQSubInstance1 = rabbitMqSub
	}

	go func() {
		count := 0

		pubsub.RabbitMQSubInstance1.ConsumeWithRetry(context.Background(), "user", "my-service-handle-for-user", "service.admin.*.create", 1, func(data string) error {
			fmt.Println(data)

			count++
			fmt.Printf("Count: %v\n", count)

			time.Sleep(1 * time.Second)

			return nil
		})
	}()

	go func() {
		for i := 1; ; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			data := fmt.Sprintf("my-payload-%v", i)
			pubsub.RabbitMQPubInstance1.PublishWithRetry(ctx, "user", fmt.Sprintf("service.admin.user-%v.create", i), data)

			cancel()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {}
}

func Example2() {
	rabbitmqclient.RabbitMQClientConnInstance = rabbitmqclient.NewRabbitMQClient(rabbitmqclient.RabbitMQConfig{
		Host:     "pitel-cx.dev.tel4vn.com",
		Port:     5672,
		Username: "pitel",
		Password: "Pitel@8229",
	})
	if err := rabbitmqclient.RabbitMQClientConnInstance.DeclareExchange("user", "topic"); err != nil {
		log.Fatalf("Failed to declare new exchange: %v", err.Error())
	}
	if rabbitMqPub, err := pubsub.NewRabbitMQPub[*model.DataStruct](); err != nil {
		log.Fatalf("Failed to create new publisher: %v", err.Error())
	} else {
		pubsub.RabbitMQPubInstance2 = rabbitMqPub
	}
	if rabbitMqSub, err := pubsub.NewRabbitMQSub[*model.DataStruct](); err != nil {
		log.Fatalf("Failed to create new subscriber: %v", err.Error())
	} else {
		pubsub.RabbitMQSubInstance2 = rabbitMqSub
	}

	go func() {
		count := 0

		pubsub.RabbitMQSubInstance2.ConsumeWithRetry(context.Background(), "user", "my-service-handle-for-user", "service.admin.*.create", 1, func(data *model.DataStruct) error {
			fmt.Println(*data)

			count++
			fmt.Printf("Count: %v\n", count)

			time.Sleep(1 * time.Second)

			return nil
		})
	}()

	go func() {
		for i := 1; ; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			data := model.DataStruct{
				Field1: fmt.Sprintf("my-payload-%v", i),
				Field2: rand.Int32(),
				Field3: rand.Int64(),
				Field4: rand.Float32(),
				Field5: rand.Float64(),
				Field6: time.Now(),
				Field7: model.SubDataStruct{
					Field1: uuid.New().String(),
					Field2: rand.Int32(),
					Field3: rand.Int64(),
				},
			}
			pubsub.RabbitMQPubInstance2.PublishWithRetry(ctx, "user", fmt.Sprintf("service.admin.user-%v.create", i), &data)

			cancel()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {}
}
