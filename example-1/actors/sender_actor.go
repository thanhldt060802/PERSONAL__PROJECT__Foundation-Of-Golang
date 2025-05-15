package actors

import (
	"fmt"
	"math/rand"
	"thanhldt060802/dto"
	"time"

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
	processName := "receiver_1"
	delayTime := time.Duration(rand.Intn(int(200*time.Millisecond-100*time.Millisecond))) + 100*time.Millisecond

	switch message.(string) {
	case "local":
		{
			process := gen.Atom(processName)

			message := dto.SimpleRequest{
				Message: fmt.Sprintf("Task %d of process %s", senderActor.count, senderActor.Name()),
			}

			senderActor.Log().Info("--> %s (LOCAL): %#v", process, message)

			result, err := senderActor.Call(process, message)
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
				Name: gen.Atom(processName),
				Node: "node2@localhost",
			}

			message := dto.SimpleRequest{
				Message: fmt.Sprintf("Task %d of process %s", senderActor.count, senderActor.Name()),
			}

			senderActor.Log().Info("--> %s (REMOTE): %#v", process.Name, message)

			result, err := senderActor.Call(process, message)
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
