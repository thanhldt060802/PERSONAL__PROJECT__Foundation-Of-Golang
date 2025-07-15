package service

import (
	"context"
	"thanhldt060802/model"
	"thanhldt060802/repository"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type (
	IUserService interface {
		Gets(ctx context.Context) ([]*model.User, error)
	}
	UserService struct {
		tracer trace.Tracer
	}
)

func NewUserService() IUserService {
	return &UserService{
		tracer: otel.Tracer("service/user.go"),
	}
}

func (s *UserService) Gets(ctx context.Context) ([]*model.User, error) {
	ctx, span1 := s.tracer.Start(ctx, "Service Gets()")
	defer span1.End()

	time.Sleep(2 * time.Second)

	data, err := repository.UserRepo.Gets(ctx)
	if err != nil {
		span1.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span1.SetAttributes(attribute.Int("total", len(data)))
	span1.SetStatus(codes.Ok, "success")
	return data, nil
}
