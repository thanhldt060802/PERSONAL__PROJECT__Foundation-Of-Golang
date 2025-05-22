package app

import (
	"thanhtldt060802/types"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverActor struct {
	act.Actor

	count int
	task  types.Task
}

func FactoryReceiverActor() gen.ProcessBehavior {
	return &ReceiverActor{}
}

func (receiverActor *ReceiverActor) Init(args ...any) error {
	if args == nil {
		receiverActor.Log().Info("started process %s %s on %s", receiverActor.PID(), receiverActor.Name(), receiverActor.Node().Name())
	} else {
		receiverActor.Log().Info("started process %s %s on %s with init value %s", receiverActor.PID(), receiverActor.Name(), receiverActor.Node().Name(), args[0].(types.Task))
		receiverActor.count = 0
		receiverActor.task = args[0].(types.Task)
	}

	return nil
}

func (receiverActor *ReceiverActor) HandleMessage(from gen.PID, message any) error {
	if message == nil {
		return gen.TerminateReasonNormal
	} else {
		switch receiverActor.task.Situation {
		case "done":
			{
				receiverActor.Log().Info("Processing task %s ... (%d)", receiverActor.task, receiverActor.count)
				receiverActor.Log().Info("Get from request: %s", message.(string))

				receiverActor.count++
				if receiverActor.count == 2 {
					return gen.TerminateReasonNormal
				} else {
					return nil
				}
			}
		case "panic":
			{
				receiverActor.Log().Info("Processing task %s ... (%d)", receiverActor.task, receiverActor.count)
				receiverActor.Log().Info("Get from request: %s", message.(string))

				receiverActor.count++
				if receiverActor.count == 2 {
					panic("Simulate crash")
				} else {
					return nil
				}
			}
		}
	}

	return nil
}

func (receiverActor *ReceiverActor) Terminate(reason error) {
	receiverActor.Log().Error("Actor terminated. Panic reason: %s", reason.Error())

	receiverActor.Log().Info("--- Update progress successful")
}
