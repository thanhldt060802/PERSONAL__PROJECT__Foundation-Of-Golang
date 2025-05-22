package app

import (
	"fmt"
	"thanhldt060802/types"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type SenderActor struct {
	act.Actor

	count int
}

func FactorySenderActor() gen.ProcessBehavior {
	return &SenderActor{}
}

func (senderActor *SenderActor) Init(args ...any) error {
	senderActor.Log().Info("started process %s %s on %s", senderActor.PID(), senderActor.Name(), senderActor.Node().Name())

	senderActor.count = 1

	return nil
}

func (senderActor *SenderActor) HandleMessage(from gen.PID, message any) error {
	switch message.(string) {
	case "idle":
		{
			sendingMessage := types.SimpleMessage{
				Message: fmt.Sprintf("Task %d of process %s", senderActor.count, senderActor.Name()),
			}

			senderActor.Log().Info("--> %s: %#v", from, sendingMessage)
			senderActor.Send(from, sendingMessage)

			senderActor.count++

			return nil
		}
	}

	senderActor.Log().Error("--- Unknown message %#v", message)
	return nil
}
