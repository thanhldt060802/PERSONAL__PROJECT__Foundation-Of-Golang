package actors

import (
	"fmt"
	"math/rand"
	"thanhldt060802/dto"
	"time"

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
	{
		request := request.(dto.SimpleRequest)
		receiverActor.Log().Info("<-- %s: %#v", from, request)

		for i := 0; i < 5; i++ {
			if rand.Intn(10) == 0 {
				panic("Simulate crash")
			}

			time.Sleep(1 * time.Second)
			i++
		}

		return fmt.Sprintf("COMPLETED %#v", request), nil
	}
}
