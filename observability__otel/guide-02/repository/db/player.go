package db

import (
	"context"
	"fmt"
	"math/rand"
	"thanhldt060802/model"
	"thanhldt060802/repository"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PlayerRepo struct {
	tracer trace.Tracer
}

func NewPlayerRepo() repository.IPlayerRepo {
	return &PlayerRepo{
		tracer: otel.Tracer("repository/db/player.go"),
	}
}

func (repo *PlayerRepo) Get(ctx context.Context) ([]*model.Player, error) {
	_, span1 := repo.tracer.Start(ctx, "Repository Get()")
	defer span1.End()

	time.Sleep(2 * time.Second)

	data := []*model.Player{}
	classes := []string{"A", "B", "C"}
	for i := 1; i <= 5; i++ {
		data = append(data, &model.Player{
			PlayerUuid: uuid.New().String(),
			Name:       fmt.Sprintf("Player %v", i),
			Class:      classes[rand.Intn(len(classes))],
			Level:      rand.Intn(100) + 1,
		})
	}

	span1.SetAttributes(attribute.Int("total", len(data)))
	span1.SetStatus(codes.Ok, "success")
	return data, nil
}
