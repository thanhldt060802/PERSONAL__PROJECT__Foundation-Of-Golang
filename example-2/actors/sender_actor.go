package actors

import (
	"context"
	"thanhldt060802/common"
	"thanhldt060802/repository"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type SenderActor struct {
	act.Actor
	receiverName string
}

func FactorySenderActor() gen.ProcessBehavior {
	return &SenderActor{}
}

func (senderActor *SenderActor) Init(args ...any) error {
	senderActor.Log().Info("STARTED PROCESS %s %s on %s", senderActor.PID(), senderActor.Name(), senderActor.Node().Name())
	senderActor.receiverName = args[0].(string)
	return nil
}

func (senderActor *SenderActor) HandleMessage(from gen.PID, message any) error {
	switch message.(string) {
	case "start":
		{
			process := gen.ProcessID{
				Name: gen.Atom(senderActor.receiverName),
				Node: "node2@localhost",
			}

			for {
				id, err := repository.TaskRepositoryInstance.GetAvailable(context.Background())
				if err != nil {
					break
				}

				message := common.TaskRequest{
					Id: id,
				}

				for {
					senderActor.Log().Info(" --> SEND REQUEST to PROCESS %s", process)
					result, err := senderActor.Call(process, message)
					if err == nil {
						senderActor.Log().Info(" <-- RECEIVED RESPONSE from PROCESS %s: %s", process, result)
						break
					}

					senderActor.Log().Warning(" --- Something wrong from PROCESS %s, retrying ... (Error: %s)", process, err.Error())
					time.Sleep(2 * time.Second)
				}
			}

			return nil
		}
	}

	senderActor.Log().Error(" --- Unknown message %#v", message)
	return nil
}
