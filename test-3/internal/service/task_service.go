package service

import (
	"context"
	"fmt"
	"thanhtldt060802/actor_model/types"
	"thanhtldt060802/internal/dto"
	"thanhtldt060802/internal/repository"
	"time"

	"ergo.services/ergo/gen"
)

type taskService struct {
	taskRepository repository.TaskRepository
	node           gen.Node
	supervisorPID  gen.PID
}

type TaskService interface {
	GetExistedReceiverNames(ctx context.Context) ([]string, error)
	SendNewTaskToReceiver(ctx context.Context, reqDTO *dto.SendNewTaskRequest) error
}

func NewTaskService(taskRepository repository.TaskRepository, node gen.Node, supervisorPID gen.PID) TaskService {
	return &taskService{
		taskRepository: taskRepository,
		node:           node,
		supervisorPID:  supervisorPID,
	}
}

func (taskService *taskService) GetExistedReceiverNames(ctx context.Context) ([]string, error) {
	existedReceiverNames := make(chan []string)
	message := types.ExistedReceiverNamesMessage{
		ReceiverNames: existedReceiverNames,
	}

	if err := taskService.node.Send(taskService.supervisorPID, message); err != nil {
		return nil, fmt.Errorf("some thing wrong on actor model: %s", err.Error())
	}

	select {
	case result := <-existedReceiverNames:
		{
			return result, nil
		}
	case <-time.After(3 * time.Second):
		{
			return nil, fmt.Errorf("time out because waiting actor model fill data to channel")
		}
	}
}

func (taskService *taskService) SendNewTaskToReceiver(ctx context.Context, reqDTO *dto.SendNewTaskRequest) error {
	foundTask, err := taskService.taskRepository.GetById(ctx, reqDTO.Body.TaskId)
	if err != nil {
		return fmt.Errorf("id of task is not valid: %s", err.Error())
	}

	message := types.NewTaskMessage{
		Receiver: reqDTO.Body.Receiver,
		TaskId:   foundTask.Id,
	}
	return taskService.node.Send(taskService.supervisorPID, message)
}
