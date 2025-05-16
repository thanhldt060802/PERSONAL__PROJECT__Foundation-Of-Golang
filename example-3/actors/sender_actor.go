package actors

import (
	"context"
	"thanhldt060802/dto"
	"thanhldt060802/repository"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type SenderActorParam struct {
	ReceiverProcessName string
}

type SenderActor struct {
	act.Actor
	SenderActorParam *SenderActorParam
}

func FactorySenderActor() gen.ProcessBehavior {
	return &SenderActor{}
}

func (senderActor *SenderActor) Init(args ...any) error {
	senderActor.Log().Info("started process %s %s on %s", senderActor.PID(), senderActor.Name(), senderActor.Node().Name())
	senderActor.SenderActorParam = args[0].(*SenderActorParam)
	return nil
}

func (senderActor *SenderActor) HandleMessage(from gen.PID, message any) error {
	switch message.(string) {
	case "start":
		{
			process := gen.ProcessID{
				Name: gen.Atom(senderActor.SenderActorParam.ReceiverProcessName),
				Node: "node2@localhost",
			}

			for {
				id, err := repository.TaskRepositoryInstance.GetAvailable(context.Background())
				if err != nil {
					break
				}

				message := dto.TaskRequest{
					Id: id,
				}

				for {
					senderActor.Log().Info("--> %s: %#v", process.Name, message)

					result, err := senderActor.Call(process, message)
					if err == nil {
						senderActor.Log().Info("<-- %s: %#v", process, result)
						break
					} else {
						senderActor.Log().Warning("--- Something wrong from PROCESS %s, retrying ... (Error: %s)", process.Name, err.Error())
						time.Sleep(1 * time.Second)
					}
				}
			}

			return nil
		}
	}

	senderActor.Log().Error("--- Unknown message %#v", message)
	return nil
}
