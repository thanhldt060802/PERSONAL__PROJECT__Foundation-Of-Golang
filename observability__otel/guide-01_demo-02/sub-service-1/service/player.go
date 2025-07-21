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
	url := fmt.Sprintf("http://localhost:8002/my-sub-service-2/v1/player/%v", playUuid)
	ctx, span, req, err := tracer.StartSpanCrossService(ctx, "GET", url)
	if err != nil {
		return nil, err
	}
	span.End()

	// Part 1 - Start

	time.Sleep(1 * time.Second)

	authHeader, _ := ctx.Value("auth-header").(string)
	req.Header.Set("Authorization", authHeader)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := http.Client{}
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
	if err := json.NewDecoder(res.Body).Decode(resWrapper); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	player := &resWrapper.Data

	// Part 1 - End

	// Part 2 - Start

	time.Sleep(1 * time.Second)

	// Part 2 - End

	span.SetStatus(codes.Ok, "success")
	return player, nil
}
