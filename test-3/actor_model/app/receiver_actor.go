package app

import (
	"context"
	"fmt"
	"math/rand"
	"thanhtldt060802/actor_model/types"
	"thanhtldt060802/internal/repository"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverActorParams struct {
	taskRepository repository.TaskRepository
	task           types.Task
}

type ReceiverActor struct {
	act.Actor

	params ReceiverActorParams
}

func FactoryReceiverActor() gen.ProcessBehavior {
	return &ReceiverActor{}
}

func (receiverActor *ReceiverActor) Init(args ...any) error {
	if args == nil {
		receiverActor.Log().Info("Started process %s %s on %s", receiverActor.PID(), receiverActor.Name(), receiverActor.Node().Name())
	} else {
		receiverActor.params = args[0].(ReceiverActorParams)
		receiverActor.Log().Info("Started process %s %s on %s with init value %s", receiverActor.PID(), receiverActor.Name(), receiverActor.Node().Name(), receiverActor.params.task)
		receiverActor.Send(receiverActor.PID(), "")
	}

	return nil
}

func (receiverActor *ReceiverActor) HandleMessage(from gen.PID, message any) error {
	if message == nil {
		return gen.TerminateReasonNormal
	} else {
		receiverActor.Log().Info("--> Start processing task task_id=%d ...", receiverActor.params.task.TaskId)

		foundTask, err := receiverActor.params.taskRepository.GetById(context.Background(), receiverActor.params.task.TaskId)
		if err != nil {
			return fmt.Errorf("id of task is not valid: %s", err.Error())
		}
		if foundTask.Status != "IN PROGRESS" {
			if err := receiverActor.params.taskRepository.Update(context.Background(), foundTask); err != nil {
				return fmt.Errorf("update task on postgresql failed: %s", err.Error())
			}
		}

		receiverActor.params.task.Progress = foundTask.Progress
		receiverActor.params.task.Target = foundTask.Target
		receiverActor.params.task.ErrorRate = foundTask.ErrorRate

		for receiverActor.params.task.Progress < receiverActor.params.task.Target {
			receiverActor.params.task.Progress++
			foundTask.Progress++
			receiverActor.Log().Info("--- Processing task task_id=%d (%d/%d)", receiverActor.params.task.TaskId, receiverActor.params.task.Progress, receiverActor.params.task.Target)
			time.Sleep(1 * time.Second)
			if rand.Intn(100) < receiverActor.params.task.ErrorRate {
				panic("Simulate crash")
			}
		}

		foundTask.Status = "COMPLETED"
		if err := receiverActor.params.taskRepository.Update(context.Background(), foundTask); err != nil {
			return fmt.Errorf("update user on postgresql failed: %s", err.Error())
		}
		receiverActor.Log().Info("--- Update task task_id=%d successful", receiverActor.params.task.TaskId)

		receiverActor.Log().Info("--> Stop processing task task_id=%d ...", receiverActor.params.task.TaskId)
		return gen.TerminateReasonNormal
	}
}

func (receiverActor *ReceiverActor) Terminate(reason error) {
	receiverActor.Log().Error("Actor terminated. Panic reason: %s", reason.Error())

	if reason.Error() != gen.TerminateReasonNormal.Error() {
		foundTask, err := receiverActor.params.taskRepository.GetById(context.Background(), receiverActor.params.task.TaskId)
		if err != nil {
			receiverActor.Log().Info("--- Update task task_id=%d failed: %s", receiverActor.params.task.TaskId, err.Error())
		}
		foundTask.Status = "IN PROGRESS"
		foundTask.Progress = receiverActor.params.task.Progress
		if err := receiverActor.params.taskRepository.Update(context.Background(), foundTask); err != nil {
			receiverActor.Log().Info("--- Update task task_id=%d on postgresql failed: %s", receiverActor.params.task.TaskId, err.Error())
		}
		receiverActor.Log().Info("--- Update task task_id=%d successful", receiverActor.params.task.TaskId)
	}
}
