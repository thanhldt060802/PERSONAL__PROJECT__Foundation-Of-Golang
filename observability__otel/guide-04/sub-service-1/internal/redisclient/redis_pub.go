package redisclient

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type RedisEnvelope struct {
	TraceContext map[string]string `json:"trace_context"`
	Payload      any               `json:"payload"`
}

type RedisPub[T any] struct {
	redisClient *redis.Client
}

type IRedisPub[T any] interface {
	Publish(ctx context.Context, channel string, data T) error
}

func NewRedisPub[T any](redisClient *redis.Client) IRedisPub[T] {
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
	if err := redisPub.redisClient.Publish(ctx, channel, payload).Err(); err != nil {
		log.Errorf("Publish %v to %v failed: %v", data, channel, err.Error())
		return err
	}

	log.Errorf("Publish %v to %v successful", data, channel)
	return nil
}
