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

type UserRepo struct {
	tracer trace.Tracer
}

func NewUserRepo() repository.IUserRepo {
	return &UserRepo{
		tracer: otel.Tracer("repository/db/user.go"),
	}
}

func (repo *UserRepo) Gets(ctx context.Context) ([]*model.User, error) {
	_, span1 := repo.tracer.Start(ctx, "Repository Gets()")
	defer span1.End()

	time.Sleep(2 * time.Second)

	data := []*model.User{}
	for i := 1; i <= 5; i++ {
		data = append(data, &model.User{
			UserUuid: uuid.New().String(),
			FullName: fmt.Sprintf("User %v", i),
			Age:      rand.Intn(30) + 1,
		})
	}

	span1.SetAttributes(attribute.Int("total", len(data)))
	span1.SetStatus(codes.Ok, "success")
	return data, nil
}
