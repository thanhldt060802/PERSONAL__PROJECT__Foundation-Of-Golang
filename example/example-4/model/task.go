package model

import (
	"github.com/uptrace/bun"
)

type Task struct {
	bun.BaseModel `bun:"table:tasks"`

	Id       int64  `bun:"id,pk,autoincrement"`
	Name     string `bun:"name,notnull"`
	Progress int64  `bun:"progress,notnull"`
	Target   int64  `bun:"target,notnull"`
	Status   string `bun:"status,notnull"`
}
