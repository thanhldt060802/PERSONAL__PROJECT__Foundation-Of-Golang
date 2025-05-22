package handler

import (
	"context"
	"net/http"
	"thanhtldt060802/internal/dto"
	"thanhtldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(api huma.API, taskService service.TaskService) *TaskHandler {
	taskHandler := &TaskHandler{
		taskService: taskService,
	}

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/get-existed-receiver-names",
		Summary:     "/get-existed-receiver-names",
		Description: "Get existed receiver names.",
		Tags:        []string{"Demo"},
	}, taskHandler.GetExistedReceiverNames)

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/send-new-task-to-receiver",
		Summary:     "/send-new-task-to-receiver",
		Description: "Send new task to receiver.",
		Tags:        []string{"Demo"},
	}, taskHandler.SendNewTaskToReceiver)

	return taskHandler
}

func (taskHandler *TaskHandler) GetExistedReceiverNames(ctx context.Context, _ *struct{}) (*dto.BodyResponse[[]string], error) {
	existedReceiverNames, err := taskHandler.taskService.GetExistedReceiverNames(ctx)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get existed receiver names failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[[]string]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get existed receiver names successful"
	res.Body.Data = existedReceiverNames
	return res, nil
}

func (taskHandler *TaskHandler) SendNewTaskToReceiver(ctx context.Context, reqDTO *dto.SendNewTaskRequest) (*dto.SuccessResponse, error) {
	if err := taskHandler.taskService.SendNewTaskToReceiver(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Send new task to receiver failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Send new task to receiver successful"
	return res, nil
}
