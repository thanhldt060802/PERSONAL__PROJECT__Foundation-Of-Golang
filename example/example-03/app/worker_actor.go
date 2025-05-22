package app

import (
	"math/rand"
	"thanhldt060802/types"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type WorkerActor struct {
	act.Actor

	taskId int64
}

func FactoryWorkerActor() gen.ProcessBehavior {
	return &WorkerActor{}
}

func (workerActor *WorkerActor) Init(args ...any) error {
	workerActor.taskId = args[0].(int64)
	workerActor.Log().Info("Started worker %s %s on %s with init task task_id=%d", workerActor.PID(), workerActor.Name(), workerActor.Node().Name(), workerActor.taskId)

	workerActor.Send(workerActor.PID(), types.DoProcessTask{})

	return nil
}

func (workerActor *WorkerActor) HandleMessage(from gen.PID, message any) error {
	switch message.(type) {
	case types.DoProcessTask:
		{
			workerActor.Log().Info("--> Start processing task task_id=%d ...", workerActor.taskId)

			for i := 1; i < 10; i++ {
				if rand.Intn(10) == 0 {
					panic("Simulate crash")
				}
				workerActor.Log().Info("--- Processing task task_id=%d ...", workerActor.taskId)
				time.Sleep(500 * time.Millisecond)
			}

			workerActor.Log().Info("--> Stop processing task task_id=%d ...", workerActor.taskId)
			return gen.TerminateReasonNormal
		}
	}

	return nil
}
