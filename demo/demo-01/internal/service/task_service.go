package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/actor_model/types"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/repository"
	"time"

	"ergo.services/ergo/gen"
)

type taskService struct {
	taskRepository repository.TaskRepository
	node           gen.Node
	supervisorPID  gen.PID
}

type TaskService interface {
	GetExistedWorkers(ctx context.Context) ([]string, error)
	RunTask(ctx context.Context, reqDTO *dto.RunTaskRequest) error
	RunTasks(ctx context.Context, reqDTO *dto.RunTasksRequest) error
}

func NewTaskService(taskRepository repository.TaskRepository, node gen.Node, supervisorPID gen.PID) TaskService {
	return &taskService{
		taskRepository: taskRepository,
		node:           node,
		supervisorPID:  supervisorPID,
	}
}

func (taskService *taskService) GetExistedWorkers(ctx context.Context) ([]string, error) {
	workerNames := make(chan []string)
	message := types.GetExistedWorkersMessage{
		WorkerNames: workerNames,
	}

	if err := taskService.node.Send(taskService.supervisorPID, message); err != nil {
		return nil, fmt.Errorf("some thing wrong on actor model: %s", err.Error())
	}

	select {
	case result := <-workerNames:
		{
			return result, nil
		}
	case <-time.After(3 * time.Second):
		{
			return nil, fmt.Errorf("time out for waiting actor model fill data to channel")
		}
	}
}

func (taskService *taskService) RunTask(ctx context.Context, reqDTO *dto.RunTaskRequest) error {
	message := types.RunTaskMessage{
		WorkerName: reqDTO.Body.WorkerName,
		TaskId:     reqDTO.Body.TaskId,
	}
	if err := taskService.node.Send(taskService.supervisorPID, message); err != nil {
		return fmt.Errorf("some thing wrong on actor model: %s", err.Error())
	}

	return nil
}

func (taskService *taskService) RunTasks(ctx context.Context, reqDTO *dto.RunTasksRequest) error {
	message := types.RunTasksMessage{
		TaskIds: reqDTO.Body.TaskIds,
	}
	if err := taskService.node.Send(taskService.supervisorPID, message); err != nil {
		return fmt.Errorf("some thing wrong on actor model: %s", err.Error())
	}

	return nil
}
