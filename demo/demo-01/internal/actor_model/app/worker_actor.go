package app

import (
	"context"
	"fmt"
	"math/rand"
	"thanhldt060802/internal/actor_model/types"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type WorkerActor struct {
	act.Actor

	taskRepository repository.TaskRepository
	taskId         int64

	task *model.Task
}

func FactoryWorkerActor() gen.ProcessBehavior {
	return &WorkerActor{}
}

func (workerActor *WorkerActor) Init(args ...any) error {
	if args == nil {
		workerActor.Log().Info("Started worker %s %s on %s", workerActor.PID(), workerActor.Name(), workerActor.Node().Name())
	} else {
		workerActor.taskRepository = args[0].(repository.TaskRepository)
		workerActor.taskId = args[1].(int64)
		workerActor.Log().Info("Started worker %s %s on %s with init task task_id=%d", workerActor.PID(), workerActor.Name(), workerActor.Node().Name(), workerActor.taskId)

		workerActor.Send(workerActor.PID(), types.DoProcessTask{})
	}

	return nil
}

func (workerActor *WorkerActor) HandleMessage(from gen.PID, message any) error {
	switch message.(type) {
	case types.DoStart:
		{
			return gen.TerminateReasonNormal
		}
	case types.DoProcessTask:
		{
			workerActor.Log().Info("--> Start processing task task_id=%d ...", workerActor.taskId)

			for {
				foundTask, err := workerActor.taskRepository.GetById(context.Background(), workerActor.taskId)
				if err != nil {
					workerActor.Log().Error("--- Get task task_id=%d from postgresql failed: %s", workerActor.taskId, err.Error())
					return fmt.Errorf("id of task is not valid: %s", err.Error())
				}
				if foundTask.Status != "IN PROGRESS" {
					workerActor.task = foundTask
					break
				}
				time.Sleep(1 * time.Second)
			}
			workerActor.Log().Info("--- Get task task_id=%d from postgresql successful", workerActor.task.Id)

			workerActor.task.Status = "IN PROGRESS"
			if err := workerActor.taskRepository.Update(context.Background(), workerActor.task); err != nil {
				workerActor.Log().Error("--- Update IN PROGRESS task task_id=%d on postgresql failed: %s", workerActor.task.Id, err.Error())
				return fmt.Errorf("update task on postgresql failed: %s", err.Error())
			}
			workerActor.Log().Info("--- Update IN PROGRESS task task_id=%d on postgresql successful", workerActor.task.Id)

			for workerActor.task.Progress < workerActor.task.Target {
				if rand.Intn(100) < workerActor.task.ErrorRate {
					panic("Simulate crash")
				}
				workerActor.task.Progress++
				workerActor.Log().Info("--- Processing task task_id=%d (%d/%d)", workerActor.task.Id, workerActor.task.Progress, workerActor.task.Target)
				time.Sleep(1 * time.Second)
			}

			workerActor.task.Status = "COMPLETED"
			if err := workerActor.taskRepository.Update(context.Background(), workerActor.task); err != nil {
				workerActor.Log().Error("--- Update COMPLETED task task_id=%d on postgresql failed: %s", workerActor.task.Id, err.Error())
				return fmt.Errorf("update task on postgresql failed: %s", err.Error())
			}
			workerActor.Log().Info("--- Update COMPLETED task task_id=%d on postgresql successful", workerActor.task.Id)

			workerActor.Log().Info("--> Stop processing task task_id=%d ...", workerActor.task.Id)
			return gen.TerminateReasonNormal
		}
	}

	return nil
}

func (workerActor *WorkerActor) Terminate(reason error) {
	if reason.Error() != gen.TerminateReasonNormal.Error() {
		workerActor.Log().Error("Actor terminated. Panic reason: %s", reason.Error())

		workerActor.task.Status = "CANCEL"
		if err := workerActor.taskRepository.Update(context.Background(), workerActor.task); err != nil {
			workerActor.Log().Error("--- Update CANCEL task task_id=%d on postgresql failed: %s", workerActor.task.Id, err.Error())
		}
		workerActor.Log().Info("--- Update CANCEL task task_id=%d on postgresql successful", workerActor.task.Id)
	}
}
