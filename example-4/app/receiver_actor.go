package app

import (
	"math/rand"
	"thanhldt060802/types"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverActorParams struct {
	dispatcherProcessName string
	dispatcherNodeName    string
}

type ReceiverActor struct {
	act.Actor

	id       int64
	progress int64
	target   int64

	params ReceiverActorParams
}

func FactoryReceiverActor() gen.ProcessBehavior {
	return &ReceiverActor{}
}

func (receiverActor *ReceiverActor) Init(args ...any) error {
	receiverActor.Log().Info("started process %s %s on %s", receiverActor.PID(), receiverActor.Name(), receiverActor.Node().Name())

	receiverActor.params = args[0].(ReceiverActorParams)

	process := gen.ProcessID{
		Name: gen.Atom(receiverActor.params.dispatcherProcessName),
		Node: gen.Atom(receiverActor.params.dispatcherNodeName),
	}
	sendingMessage := types.ReturnTaskMessage{
		SelfStatus: "idle",
	}
	receiverActor.Send(process, sendingMessage)

	return nil
}

func (receiverActor *ReceiverActor) HandleMessage(from gen.PID, message any) error {
	receiverMessage := message.(types.DoTaskMessage)

	receiverActor.Log().Info("<-- %s: %#v", from, receiverMessage)

	receiverActor.id = receiverMessage.Id
	receiverActor.progress = receiverMessage.Progress
	receiverActor.target = receiverMessage.Target

	defer func() {
		if r := recover(); r != nil {
			process := gen.ProcessID{
				Name: gen.Atom(receiverActor.params.dispatcherProcessName),
				Node: gen.Atom(receiverActor.params.dispatcherNodeName),
			}
			sendingMessage := types.ReturnTaskMessage{
				Id:         receiverActor.id,
				Progress:   receiverActor.progress,
				Target:     receiverActor.target,
				SelfStatus: "cancel",
			}
			if _, err := receiverActor.Call(process, sendingMessage); err != nil {
				receiverActor.Log().Error("Failed to notify dispatcher before crash: %s", err.Error())
			}

			panic(r)
		}
	}()

	for receiverActor.progress < receiverActor.target {
		if rand.Intn(10) == 0 {
			panic("Simulate crash")
		}

		time.Sleep(1 * time.Second)
		receiverActor.progress++
	}

	receiverActor.Log().Info("COMPLETED %#v", message)

	sendingMessage1 := types.ReturnTaskMessage{
		Id:         receiverActor.id,
		Progress:   receiverActor.progress,
		Target:     receiverActor.target,
		SelfStatus: "complete",
	}
	receiverActor.Send(from, sendingMessage1)

	sendingMessage2 := types.ReturnTaskMessage{
		SelfStatus: "idle",
	}
	receiverActor.Send(from, sendingMessage2)

	return nil
}

func (receiverActor *ReceiverActor) Terminate(reason error) {
	receiverActor.Log().Error("Actor terminated. Panic reason: %s", reason.Error())
}
