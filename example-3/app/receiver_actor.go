package app

import (
	"math/rand"
	"thanhldt060802/types"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverActorParam struct {
	SenderName     string
	SenderNodeName string
}

type ReceiverActor struct {
	act.Actor

	params ReceiverActorParam
}

func FactoryReceiverActor() gen.ProcessBehavior {
	return &ReceiverActor{}
}

func (receiverActor *ReceiverActor) Init(args ...any) error {
	receiverActor.Log().Info("started process %s %s on %s", receiverActor.PID(), receiverActor.Name(), receiverActor.Node().Name())

	receiverActor.params = args[0].(ReceiverActorParam)

	process := gen.ProcessID{
		Name: gen.Atom(receiverActor.params.SenderName),
		Node: gen.Atom(receiverActor.params.SenderNodeName),
	}
	receiverActor.SendAfter(process, "idle", 2*time.Second)

	return nil
}

func (receiverActor *ReceiverActor) HandleMessage(from gen.PID, message any) error {
	receivedMessage := message.(types.SimpleMessage)

	receiverActor.Log().Info("<-- %s: %#v", from, receivedMessage)

	for i := 0; i < 5; i++ {
		if rand.Intn(10) == 0 {
			panic("Simulate crash")
		}

		time.Sleep(1 * time.Second)
		i++
	}

	receiverActor.Log().Info("--- Processing successful: %#v", receivedMessage)
	receiverActor.Send(from, "idle")

	return nil
}

func (receiverActor *ReceiverActor) Terminate(reason error) {
	receiverActor.Log().Error("Actor terminated. Panic reason: %s", reason.Error())

	receiverActor.Log().Info("--- Update progress successful")
}
