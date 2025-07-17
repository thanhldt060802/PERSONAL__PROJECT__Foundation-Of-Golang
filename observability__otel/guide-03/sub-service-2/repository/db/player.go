package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"thanhldt060802/common/tracer"
	"thanhldt060802/internal/sqlclient"
	"thanhldt060802/model"
	"thanhldt060802/repository"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/codes"
)

type PlayerRepo struct {
}

func NewPlayerRepo() repository.IPlayerRepo {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repo := &PlayerRepo{}
	repo.DeleteTable(ctx)
	repo.InitTable(ctx)
	repo.GenerateData(ctx)

	return repo
}

func (repo *PlayerRepo) DeleteTable(ctx context.Context) {
	if err := repository.DropTable(sqlclient.BunSqlClientConnInstance, ctx, (*model.Player)(nil)); err != nil {
		panic(err)
	}
}

func (repo *PlayerRepo) InitTable(ctx context.Context) {
	if err := repository.CreateTable(sqlclient.BunSqlClientConnInstance, ctx, (*model.Player)(nil)); err != nil {
		panic(err)
	}
}

func (repo *PlayerRepo) GenerateData(ctx context.Context) {
	classes := []string{"Assassin", "Warrior", "Mage", "Gunner"}

	if err := sqlclient.BunSqlClientConnInstance.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for i := 1; i <= 30; i++ {
			user := model.Player{
				PlayerUuid: uuid.New().String(),
				Name:       fmt.Sprintf("Player %v", i),
				Class:      classes[rand.Intn(len(classes))],
				Level:      (rand.Intn(10) + 1) * 10,
			}
			if _, err := tx.NewInsert().Model(&user).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func (repo *PlayerRepo) GetById(ctx context.Context, playUuid string) (*model.Player, error) {
	ctx, span := tracer.StartSpanInternal(ctx, "repository/db/player.go", "Repository.GetById")
	defer span.End()

	time.Sleep(1 * time.Second)

	player := new(model.Player)

	query := sqlclient.BunSqlClientConnInstance.GetDB().NewSelect().Model(player).
		Where("player_uuid = ?", playUuid)

	err := query.Scan(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	} else {
		span.SetStatus(codes.Ok, "success")
		return player, nil
	}
}
