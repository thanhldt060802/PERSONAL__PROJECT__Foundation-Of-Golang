package service

import (
	"context"
	"encoding/json"
	"thanhldt060802/common/tracer"
	"thanhldt060802/internal/redisclient"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type (
	IPlayerService interface {
		GetById(ctx context.Context, playUuid string) (string, error)
	}
	PlayerService struct {
	}
)

func NewPlayerService() IPlayerService {
	return &PlayerService{}
}

func (s *PlayerService) GetById(ctx context.Context, playUuid string) (string, error) {
	_, span1 := tracer.StartSpan(ctx, "service/player.go", "Service.GetById")
	defer span1.End()

	// Part 1 - Start

	time.Sleep(1 * time.Second)

	pubCtx, span2 := tracer.StartSpan(ctx, "service/player.go", "RedisPubSub.Publish")

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(pubCtx, carrier)

	envelope := redisclient.RedisEnvelope{
		TraceContext: carrier,
		Payload:      playUuid,
	}

	redisPub := redisclient.NewRedisPub[*redisclient.RedisEnvelope](redisclient.RedisClient.GetClient())

	span2.SetAttributes(attribute.String("redis.pub_channel", "test.trace.pubsub"))
	if err := redisPub.Publish(ctx, "test.trace.pubsub", &envelope); err != nil {
		span2.RecordError(err)
		span2.SetStatus(codes.Error, err.Error())
		span2.End()
		return "", err
	}
	payloadBytes, _ := json.Marshal(envelope.Payload)
	span2.SetAttributes(attribute.String("test.trace.pubsub.payload", string(payloadBytes)))

	span2.SetStatus(codes.Ok, "success")
	span2.End()

	// Part 1 - End

	// Part 2 - Start

	time.Sleep(1 * time.Second)

	span1.SetAttributes(attribute.String("other.action", "action.result"))

	// Part 2 - End

	span1.SetStatus(codes.Ok, "success")
	return "success", nil
}
