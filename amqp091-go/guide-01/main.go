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

		if err := pubsub.RabbitMQSubInstance1.Consume(context.Background(), "user", "user-handler-queue", "service.admin.*.create", func(data string) error {
			fmt.Println(data)

			count++
			fmt.Printf("Count: %v\n", count)

			time.Sleep(1 * time.Second)

			return nil
		}); err != nil {
			log.Fatalf("Failed to start consume on %v for %v of %v: %v", "user-handler-queue", "service.admin.*.create", "user", err.Error())
		}
	}()

	go func() {
		for i := 1; ; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			data := fmt.Sprintf("my-payload-%v", i)
			if err := pubsub.RabbitMQPubInstance1.Publish(ctx, "user", fmt.Sprintf("service.admin.user-%v.create", i), data); err != nil {
				cancel()
				return
			}

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

		if err := pubsub.RabbitMQSubInstance2.Consume(context.Background(), "user", "user-handler-queue", "service.admin.*.create", func(data *model.DataStruct) error {
			fmt.Println(*data)

			count++
			fmt.Printf("Count: %v\n", count)

			time.Sleep(1 * time.Second)

			return nil
		}); err != nil {
			log.Fatalf("Failed to start consume on %v for %v of %v: %v", "user-handler-queue", "service.admin.*.create", "user", err.Error())
		}
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
			if err := pubsub.RabbitMQPubInstance2.Publish(ctx, "user", fmt.Sprintf("service.admin.user-%v.create", i), &data); err != nil {
				cancel()
				return
			}

			cancel()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {}
}
