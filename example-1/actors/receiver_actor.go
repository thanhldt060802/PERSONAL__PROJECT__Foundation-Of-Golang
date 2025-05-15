package actors

import (
	"fmt"
	"thanhldt060802/common"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverActor struct {
	act.Actor
}

func FactoryReceiverActor() gen.ProcessBehavior {
	return &ReceiverActor{}
}

func (receiverActor *ReceiverActor) Init(args ...any) error {
	receiverActor.Log().Info("STARTED PROCESS %s %s on %s", receiverActor.PID(), receiverActor.Name(), receiverActor.Node().Name())
	return nil
}

func (receiverActor *ReceiverActor) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	switch request := request.(type) {
	case common.LocalRequest:
		{
			receiverActor.Log().Info(" <-- RECEIVED LOCAL REQUEST from %s: %s", from, request.Message)
			return fmt.Sprintf("DONE %s", request.Message), nil
		}
	case common.RemoteRequest:
		{
			receiverActor.Log().Info(" <-- RECEIVED REMOTE REQUEST from %s: %s", from, request.Message)
			return fmt.Sprintf("DONE %s", request.Message), nil
		}
	}

	receiverActor.Log().Info("Received unknown request: %#v", request)
	return nil, nil
}
