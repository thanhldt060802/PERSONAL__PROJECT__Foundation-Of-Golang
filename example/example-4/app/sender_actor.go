package app

import (
	"context"
	"thanhldt060802/repository"
	"thanhldt060802/types"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type SenderActor struct {
	act.Actor
}

func FactorySenderActor() gen.ProcessBehavior {
	return &SenderActor{}
}

func (senderActor *SenderActor) Init(args ...any) error {
	senderActor.Log().Info("started process %s %s on %s", senderActor.PID(), senderActor.Name(), senderActor.Node().Name())
	return nil
}

func (senderActor *SenderActor) HandleMessage(from gen.PID, message any) error {
	switch message.(string) {
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

			senderActor.Log().Info("--> %s: %#v", from, sendingMessage)
			senderActor.Send(from, sendingMessage)

			return nil
		}
	}

	senderActor.Log().Error("--- Unknown message %#v", message)
	return nil
}
