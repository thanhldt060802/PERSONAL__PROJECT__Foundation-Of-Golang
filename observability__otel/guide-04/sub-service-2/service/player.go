package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"thanhldt060802/common/tracer"
	"thanhldt060802/internal/redisclient"
	"thanhldt060802/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type (
	IPlayerService interface {
		InitSubscriber()
	}
	PlayerService struct {
	}
)

func NewPlayerService() IPlayerService {
	return &PlayerService{}
}

func (s *PlayerService) InitSubscriber() {
	redisSub := redisclient.NewRedisSub[*redisclient.RedisEnvelope](redisclient.RedisClient.GetClient())

	redisSub.Subscribe(context.Background(), "test.trace.pubsub", func(data *redisclient.RedisEnvelope) {
		carrier := propagation.MapCarrier(data.TraceContext)
		pubCtx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

		subCtx, span := tracer.StartSpan(pubCtx, "service/player.go", "RedisPubSub.Subscribe")
		defer span.End()

		span.SetAttributes(attribute.String("redis.sub_channel", "test.trace.pubsub"))
		payloadBytes, _ := json.Marshal(data.Payload)
		span.SetAttributes(attribute.String("test.trace.pubsub.payload", string(payloadBytes)))

		playerUuid, ok := data.Payload.(string)
		if !ok {
			err := errors.New("invalid payload")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}

		player, err := repository.PlayerRepo.GetById(subCtx, playerUuid)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		fmt.Println(*player)

		span.SetStatus(codes.Ok, "success")
	})
}
