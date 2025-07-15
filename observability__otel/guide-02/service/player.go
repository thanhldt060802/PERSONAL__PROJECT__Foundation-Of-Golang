package service

import (
	"context"
	"thanhldt060802/common/tracer"
	"thanhldt060802/model"
	"thanhldt060802/repository"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type (
	IPlayerService interface {
		Get(ctx context.Context) ([]*model.Player, error)
	}
	PlayerService struct {
	}
)

func NewPlayerService() IPlayerService {
	return &PlayerService{}
}

func (s *PlayerService) Get(ctx context.Context) ([]*model.Player, error) {
	ctx, span1 := tracer.TracingMg.NewSpan(ctx, "service/player.go", "Service -> Get()")
	defer span1.End()

	time.Sleep(2 * time.Second)

	data, err := repository.PlayerRepo.Get(ctx)
	if err != nil {
		span1.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span1.SetAttributes(attribute.Int("total", len(data)))
	span1.SetStatus(codes.Ok, "success")
	return data, nil
}
