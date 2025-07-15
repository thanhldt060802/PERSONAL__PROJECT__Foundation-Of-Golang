package repository

import (
	"context"
	"thanhldt060802/model"
)

type IUserRepo interface {
	Gets(ctx context.Context) ([]*model.User, error)
}

var UserRepo IUserRepo
