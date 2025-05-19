package app

import (
	"context"
	"thanhldt060802/repository"
	"thanhldt060802/types"

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
	receivedMessage := message.(string)

	switch receivedMessage {
	case "idle":
		{
			task, err := repository.TaskRepositoryInstance.GetAvailable(context.Background())
			if err != nil {
				return nil
			}

			sendingMessage := types.DoTaskMessage{
				TaskId:       task.Id,
				TaskProgress: task.Progress,
				TaskTarget:   task.Target,
			}

			dispatcherActor.Log().Info("--> %s: %#v", from, task)
			dispatcherActor.Send(from, sendingMessage)

			return nil
		}
	}

	dispatcherActor.Log().Error("--- Unknown message %#v", message)
	return nil
}
