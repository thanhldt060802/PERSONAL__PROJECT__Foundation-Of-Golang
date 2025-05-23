package app

import (
	"thanhldt060802/types"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type WorkerActor struct {
	act.Actor

	taskDuration int
}

func FactoryWorkerActor() gen.ProcessBehavior {
	return &WorkerActor{}
}

func (workerActor *WorkerActor) Init(args ...any) error {
	workerActor.taskDuration = args[0].(int)
	workerActor.Log().Info("Started worker %s %s on %s with init task task_id=%d", workerActor.PID(), workerActor.Name(), workerActor.Node().Name())

	return nil
}

func (workerActor *WorkerActor) HandleMessage(from gen.PID, message any) error {
	switch message.(type) {
	case types.DoProcessTask:
		{
			workerActor.Log().Info("--> Start processing task")

			workerActor.Log().Info("--- Processing task")
			time.Sleep(time.Second * time.Duration(workerActor.taskDuration))

			workerActor.Log().Info("--> Stop processing task")
		}
	}

	return nil
}
