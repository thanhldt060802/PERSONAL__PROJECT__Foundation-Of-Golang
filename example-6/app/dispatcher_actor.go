package app

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type DispatcherActor struct {
	act.Actor
}

func FactoryDispatcherActor() gen.ProcessBehavior {
	return &DispatcherActor{}
}

func (dispatcherActor *DispatcherActor) Init(args ...any) error {
	dispatcherActor.Log().Info("started process %s %s on %s", dispatcherActor.PID(), dispatcherActor.Name(), dispatcherActor.Node().Name())
	return nil
}

func (dispatcherActor *DispatcherActor) HandleMessage(from gen.PID, message any) error {
	return nil
}
