package actors

import (
	"fmt"
	"math/rand"
	"thanhldt060802/common"
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
	senderActor.Log().Info("STARTED PROCESS %s %s on %s", senderActor.PID(), senderActor.Name(), senderActor.Node().Name())
	senderActor.count = 1
	return nil
}

func (senderActor *SenderActor) HandleMessage(from gen.PID, message any) error {
	processName := "receiver-1"
	delayTime := time.Duration(rand.Intn(int(200*time.Millisecond-100*time.Millisecond))) + 100*time.Millisecond

	switch message.(type) {
	case common.DoCallLocal:
		{
			process := gen.Atom(processName)

			message := common.RemoteRequest{
				Message: fmt.Sprintf("Task %d of process %s %s", senderActor.count, senderActor.PID(), senderActor.Name()),
			}

			senderActor.Log().Info(" --> SEND REQUEST to LOCAL PROCESS %s: %s", process, message.Message)

			result, err := senderActor.Call(process, message)
			if err == nil {
				senderActor.Log().Info(" <-- RECEIVED RESPONSE from LOCAL PROCESS %s: %s", process, result)
			} else {
				senderActor.Log().Error(" --- SEND REQUEST to LOCAL PROCESS %s failed: %s", process, err.Error())
			}

			senderActor.count++
			senderActor.SendAfter(senderActor.PID(), common.DoCallRemote{}, delayTime)
			return nil
		}
	case common.DoCallRemote:
		{
			process := gen.ProcessID{
				Name: gen.Atom(processName),
				Node: "node2@localhost",
			}

			message := common.RemoteRequest{
				Message: fmt.Sprintf("Task %d of process %s %s", senderActor.count, senderActor.PID(), senderActor.Name()),
			}

			senderActor.Log().Info(" --> SEND REQUEST to REMOTE PROCESS %s: %s", process.Name, message.Message)

			result, err := senderActor.Call(process, message)
			if err == nil {
				senderActor.Log().Info(" <-- RECEIVED RESPONSE from REMOTE PROCESS %s: %s", process.Name, result)
			} else {
				senderActor.Log().Error(" --- SEND REQUEST to REMOTE PROCESS %s failed: %s", process.Name, err.Error())
			}

			senderActor.count++
			senderActor.SendAfter(senderActor.PID(), common.DoCallLocal{}, delayTime)
			return nil
		}
	}

	senderActor.Log().Error(" --- Unknown message %#v", message)
	return nil
}
