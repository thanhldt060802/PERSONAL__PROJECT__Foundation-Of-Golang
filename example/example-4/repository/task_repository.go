package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/model"
)

var TaskRepositoryInstance *taskRepository

type taskRepository struct {
}

type TaskRepository interface {
	GetAvailable(ctx context.Context) (*model.Task, error)
	GetById(ctx context.Context, id int64) (*model.Task, error)
	Update(ctx context.Context, updatedTask *model.Task) error
}

func InitTaskRepository() {
	TaskRepositoryInstance = &taskRepository{}
}

func (taskRepository *taskRepository) GetAvailable(ctx context.Context) (*model.Task, error) {
	task := &model.Task{}

	if err := infrastructure.PostgresDB.NewSelect().Model(task).Where("status = 'CANCEL'").Limit(1).Scan(ctx); err != nil {
		if err := infrastructure.PostgresDB.NewSelect().Model(task).Where("status = 'PENDING'").Limit(1).Scan(ctx); err != nil {
			return nil, err
		}
	}

	task.Status = "IN PROGRESS"
	if _, err := infrastructure.PostgresDB.NewUpdate().Model(task).Where("id = ?", task.Id).Exec(ctx); err != nil {
		return nil, err
	}

	return task, nil
}

func (taskRepository *taskRepository) GetById(ctx context.Context, id int64) (*model.Task, error) {
	var task model.Task

	if err := infrastructure.PostgresDB.NewSelect().Model(&task).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return &task, nil
}

func (taskRepository *taskRepository) Update(ctx context.Context, updatedTask *model.Task) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedTask).Where("id = ?", updatedTask.Id).Exec(ctx)
	return err
}
