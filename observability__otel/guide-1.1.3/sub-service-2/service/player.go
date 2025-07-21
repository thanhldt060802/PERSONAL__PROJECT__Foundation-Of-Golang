package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"thanhldt060802/common/pubsub"
	"thanhldt060802/common/tracer"
	"thanhldt060802/repository"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	pubsub.RedisSubInstance.Subscribe(context.Background(), "test.trace.pubsub", func(data *tracer.MessageTracing) {
		pubCtx := data.ExtractTraceContext()
		subCtx, span := tracer.StartSpanInternal(pubCtx, "service/player.go", "RedisPubSub.Subscribe")
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
