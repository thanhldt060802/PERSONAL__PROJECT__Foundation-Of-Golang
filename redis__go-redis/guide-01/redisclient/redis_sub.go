package redisclient

import (
	"context"
	"encoding/json"
	"reflect"

	log "github.com/sirupsen/logrus"
)

type RedisSub[T any] struct {
	redisClient *RedisClient
}

type IRedisSub[T any] interface {
	Subscribe(ctx context.Context, channel string, handler func(data T))
}

func NewRedisSub[T any](redisClient *RedisClient) IRedisSub[T] {
	return &RedisSub[T]{
		redisClient: redisClient,
	}
}

func (redisSub *RedisSub[T]) Subscribe(ctx context.Context, channel string, handler func(data T)) {
	sub := redisSub.redisClient.redisClient.Subscribe(ctx, channel)
	ch := sub.Channel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				sub.Close()
				return
			case message := <-ch:
				var value T
				t := reflect.TypeOf(value)

				var instance any
				if t.Kind() == reflect.Ptr {
					// T is pointer to struct: create *Struct
					instance = reflect.New(t.Elem()).Interface()
				} else {
					// T is value: create pointer to value (e.g., *int, *string)
					instance = reflect.New(t).Interface()
				}

				if err := json.Unmarshal([]byte(message.Payload), instance); err != nil {
					log.Errorf("Unmarshal %v failed: %v", message.Payload, err.Error())
					continue
				}

				var data T
				if t.Kind() == reflect.Ptr {
					// T is pointer already
					data = instance.(T)
				} else {
					// T is value, dereference pointer
					data = reflect.ValueOf(instance).Elem().Interface().(T)
				}

				handler(data)
			}
		}
	}()
}
