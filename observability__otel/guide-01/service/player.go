package service

import (
	"context"
	"thanhldt060802/common/tracer"
	"thanhldt060802/model"
	"thanhldt060802/repository"
	"time"
)

type (
	IPlayerService interface {
		GetById(ctx context.Context, playUuid string) (*model.Player, error)
	}
	PlayerService struct {
	}
)

func NewPlayerService() IPlayerService {
	return &PlayerService{}
}

func (s *PlayerService) GetById(ctx context.Context, playUuid string) (*model.Player, error) {
	ctx, span := tracer.StartSpanInternal(ctx)
	defer span.End()

	// Part 1 - Start

	time.Sleep(1 * time.Second)

	player, err := repository.PlayerRepo.GetById(ctx, playUuid)
	if err != nil {
		span.Err = err
		return nil, err
	}

	// Part 1 - End

	// Part 2 - Start

	time.Sleep(1 * time.Second)

	// Part 2 - End

	return player, nil
}
