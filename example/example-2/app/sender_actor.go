package app

import (
	"fmt"
	"math/rand"
	"thanhldt060802/types"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type SenderActorParams struct {
	ReceiverName     string
	ReceiverNodeName string
}

type SenderActor struct {
	act.Actor

	params SenderActorParams

	count int
}

func FactorySenderActor() gen.ProcessBehavior {
	return &SenderActor{}
}

func (senderActor *SenderActor) Init(args ...any) error {
	senderActor.Log().Info("started process %s %s on %s", senderActor.PID(), senderActor.Name(), senderActor.Node().Name())

	senderActor.params = args[0].(SenderActorParams)
	senderActor.count = 1

	senderActor.SendAfter(senderActor.PID(), "local", 2*time.Second)

	return nil
}

func (senderActor *SenderActor) HandleMessage(from gen.PID, message any) error {
	delayTime := time.Duration(rand.Intn(int(200*time.Millisecond-100*time.Millisecond))) + 100*time.Millisecond

	sendingRequest := types.SimpleRequest{
		Message: fmt.Sprintf("Task %d of process %s", senderActor.count, senderActor.Name()),
	}

	switch message.(string) {
	case "local":
		{
			process := gen.Atom(senderActor.params.ReceiverName)

			senderActor.Log().Info("--> %s (LOCAL): %#v", process, sendingRequest)
			result, err := senderActor.Call(process, sendingRequest)
			if err == nil {
				senderActor.Log().Info("<-- %s (LOCAL): %#v", process, result)
			} else {
				senderActor.Log().Error("--- Call to %s (LOCAL) failed: %s", process, err.Error())
			}

			senderActor.count++
			senderActor.SendAfter(senderActor.PID(), "remote", delayTime)
			return nil
		}
	case "remote":
		{
			process := gen.ProcessID{
				Name: gen.Atom(senderActor.params.ReceiverName),
				Node: gen.Atom(senderActor.params.ReceiverNodeName),
			}

			senderActor.Log().Info("--> %s (REMOTE): %#v", process.Name, sendingRequest)

			result, err := senderActor.Call(process, sendingRequest)
			if err == nil {
				senderActor.Log().Info("<-- %s (REMOTE): %#v", process.Name, result)
			} else {
				senderActor.Log().Error("--- Call to %s (REMOTE) failed: %s", process.Name, err.Error())
			}

			senderActor.count++
			senderActor.SendAfter(senderActor.PID(), "local", delayTime)
			return nil
		}
	}

	senderActor.Log().Error("--- Unknown message %#v", message)
	return nil
}
