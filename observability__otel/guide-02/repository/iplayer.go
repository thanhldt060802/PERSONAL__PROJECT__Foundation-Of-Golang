package repository

import (
	"context"
	"thanhldt060802/model"
)

type IPlayerRepo interface {
	Get(ctx context.Context) ([]*model.Player, error)
}

var PlayerRepo IPlayerRepo
