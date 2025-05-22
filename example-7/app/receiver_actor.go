package app

import (
	"fmt"
	"thanhldt060802/types"

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
	receiverActor.Log().Info("started process %s %s on %s", receiverActor.PID(), receiverActor.Name(), receiverActor.Node().Name())
	return nil
}

func (receiverActor *ReceiverActor) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	receiverRequest := request.(types.SimpleRequest)
	receiverActor.Log().Info("<-- %s: %#v", from, receiverRequest)
	receiverActor.Log().Info("--- Processing successful: %#v", receiverRequest)
	return fmt.Sprintf("COMPLETED %#v", receiverRequest), nil
}
