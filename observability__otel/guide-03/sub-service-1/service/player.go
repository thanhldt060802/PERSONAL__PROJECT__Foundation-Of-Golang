package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"thanhldt060802/common/tracer"
	"thanhldt060802/model"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
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
	ctx, span := tracer.StartSpan(ctx, "service/player.go", "Service.GetById")
	defer span.End()

	// Part 1 - Start

	time.Sleep(1 * time.Second)

	client := http.Client{}
	url := fmt.Sprintf("http://localhost:8002/my-sub-service-2/v1/player/%v", playUuid)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	authHeader, _ := ctx.Value("auth-header").(string)
	req.Header.Set("Authorization", authHeader)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err := errors.New("response failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	resWrapper := new(struct {
		Data model.Player
	})
	err = json.NewDecoder(res.Body).Decode(resWrapper)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	player := &resWrapper.Data
	span.SetAttributes(attribute.String("player.player_uuid", player.PlayerUuid))

	// Part 1 - End

	// Part 2 - Start

	time.Sleep(1 * time.Second)

	span.SetAttributes(attribute.String("other.action", "action.result"))

	// Part 2 - End

	span.SetStatus(codes.Ok, "success")
	return player, nil
}
