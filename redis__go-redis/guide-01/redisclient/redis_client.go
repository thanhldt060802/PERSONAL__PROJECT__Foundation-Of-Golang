package redisclient

import (
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	redisClient *redis.Client
}

func NewRedisClient(redisOption *redis.Options) *RedisClient {
	return &RedisClient{
		redisClient: redis.NewClient(redisOption),
	}
}
