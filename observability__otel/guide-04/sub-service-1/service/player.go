package service

import (
	"context"
	"encoding/json"
	"thanhldt060802/common/pubsub"
	"thanhldt060802/common/tracer"
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
	_, span1 := tracer.StartSpanInternal(ctx, "service/player.go", "Service.GetById")
	defer span1.End()

	// Part 1 - Start

	time.Sleep(1 * time.Second)

	pubCtx, span2 := tracer.StartSpanInternal(ctx, "service/player.go", "RedisPubSub.Publish")

	msgTrace := tracer.MessageTracing{
		TraceContext: propagation.MapCarrier{},
		Payload:      playUuid,
	}

	otel.GetTextMapPropagator().Inject(pubCtx, msgTrace.TraceContext)

	span2.SetAttributes(attribute.String("redis.pub_channel", "test.trace.pubsub"))
	if err := pubsub.RedisPubInstance.Publish(ctx, "test.trace.pubsub", &msgTrace); err != nil {
		span2.RecordError(err)
		span2.SetStatus(codes.Error, err.Error())
		span2.End()
		return "", err
	}
	payloadBytes, _ := json.Marshal(msgTrace.Payload)
	span2.SetAttributes(attribute.String("test.trace.pubsub.payload", string(payloadBytes)))

	span2.SetStatus(codes.Ok, "success")
	span2.End()

	// Part 1 - End

	// Part 2 - Start

	time.Sleep(1 * time.Second)

	// Part 2 - End

	span1.SetStatus(codes.Ok, "success")
	return "success", nil
}
