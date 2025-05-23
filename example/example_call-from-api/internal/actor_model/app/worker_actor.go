package app

import (
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

			workerActor.task = workerActor.taskRepository.GetById(workerActor.taskId)
			workerActor.Log().Info("--- Get task task_id=%d from database successful", workerActor.task.Id)

			workerActor.Log().Info("--- Update IN PROGRESS task task_id=%d on database successful", workerActor.task.Id)

			for workerActor.task.Progress < workerActor.task.Target {
				if rand.Intn(10) == 0 {
					panic("Simulate crash")
				}
				workerActor.task.Progress++
				workerActor.Log().Info("--- Processing task task_id=%d (%d/%d)", workerActor.task.Id, workerActor.task.Progress, workerActor.task.Target)
				time.Sleep(500 * time.Millisecond)
			}

			workerActor.Log().Info("--- Update COMPLETED task task_id=%d on database successful", workerActor.task.Id)

			workerActor.Log().Info("--> Stop processing task task_id=%d ...", workerActor.task.Id)
			return gen.TerminateReasonNormal
		}
	}

	return nil
}

func (workerActor *WorkerActor) Terminate(reason error) {
	if reason.Error() != gen.TerminateReasonNormal.Error() {
		workerActor.Log().Error("Actor terminated. Panic reason: %s", reason.Error())

		workerActor.Log().Info("--- Update CANCEL task task_id=%d on database successful", workerActor.task.Id)
	}
}
