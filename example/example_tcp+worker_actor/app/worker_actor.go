package app

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type WorkerActor struct {
	act.Actor
}

func FactoryWorkerActor() gen.ProcessBehavior {
	return &WorkerActor{}
}

func (workerActor *WorkerActor) Init(args ...any) error {
	workerActor.Log().Info("Started process successful")
	return nil
}

func (workerActor *WorkerActor) HandleMessage(from gen.PID, message any) error {
	workerActor.Log().Info("--> Start processing task")

	workerActor.Log().Info("--- Processing task")

	workerActor.Log().Info("--> Stop processing task")
	return nil
}
