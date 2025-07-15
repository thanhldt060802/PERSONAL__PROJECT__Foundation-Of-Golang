package redisclient

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type RedisPub[T any] struct {
	redisClient *RedisClient
}

type IRedisPub[T any] interface {
	Publish(ctx context.Context, channel string, data T) error
}

func NewRedisPub[T any](redisClient *RedisClient) IRedisPub[T] {
	return &RedisPub[T]{
		redisClient: redisClient,
	}
}

func (redisPub *RedisPub[T]) Publish(ctx context.Context, channel string, data T) error {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Marshal data failed: %v", err.Error())
		return err
	}
	if err := redisPub.redisClient.redisClient.Publish(ctx, channel, payload).Err(); err != nil {
		log.Errorf("Publish %v to %v failed: %v", data, channel, err.Error())
		return err
	}

	log.Errorf("Publish %v to %v successful", data, channel)
	return nil
}
