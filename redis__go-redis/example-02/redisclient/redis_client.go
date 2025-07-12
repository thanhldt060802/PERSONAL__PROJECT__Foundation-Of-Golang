package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	redisClient *redis.Client
}

type IRedisClient interface {
	Publish(ctx context.Context, channel string, data any) error
	Subscribe(ctx context.Context, channel string, handler func(payload string))
}

func NewRedisClient(redisOption *redis.Options) IRedisClient {
	return &RedisClient{
		redisClient: redis.NewClient(redisOption),
	}
}

func (redisClient *RedisClient) Publish(ctx context.Context, channel string, message any) error {
	return redisClient.redisClient.Publish(ctx, channel, message).Err()
}

func (redisClient *RedisClient) Subscribe(ctx context.Context, channel string, handler func(payload string)) {
	sub := redisClient.redisClient.Subscribe(ctx, channel)
	ch := sub.Channel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				sub.Close()
				return
			case message := <-ch:
				handler(message.Payload)
			}
		}
	}()
}
