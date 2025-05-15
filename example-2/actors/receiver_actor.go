package actors

import (
	"context"
	"fmt"
	"math/rand"
	"thanhldt060802/common"
	"thanhldt060802/model"
	"thanhldt060802/repository"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverActor struct {
	act.Actor
	Task *model.Task
}

func FactoryReceiverActor() gen.ProcessBehavior {
	return &ReceiverActor{}
}

func (receiverActor *ReceiverActor) Init(args ...any) error {
	receiverActor.Log().Info("STARTED PROCESS %s %s on %s", receiverActor.PID(), receiverActor.Name(), receiverActor.Node().Name())
	return nil
}

func (receiverActor *ReceiverActor) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	{
		request := request.(common.TaskRequest)
		receiverActor.Log().Info(" <-- RECEIVED REQUEST from %s", from)

		foundTask, err := repository.TaskRepositoryInstance.GetById(context.Background(), request.Id)
		if err != nil {
			return nil, fmt.Errorf("id of task is not valid")
		}

		receiverActor.Task = foundTask

		for receiverActor.Task.Progress < receiverActor.Task.Target {
			number := rand.Intn(7)
			if number == 0 {
				panic("Simulate crash")
			}

			time.Sleep(1 * time.Second)
			receiverActor.Task.Progress++
		}
		receiverActor.Task.Status = "DONE"

		if err := repository.TaskRepositoryInstance.Update(context.Background(), receiverActor.Task); err != nil {
			return nil, fmt.Errorf("update task failed")
		}

		return fmt.Sprintf("DONE %#v", request), nil
	}
}

func (receiverActor *ReceiverActor) Terminate(reason error) {
	receiverActor.Log().Error("Actor terminated. Panic reason: %s", reason.Error())

	if err := repository.TaskRepositoryInstance.Update(context.Background(), receiverActor.Task); err != nil {
		receiverActor.Log().Error("Update task failed: %s", err.Error())
	}
}
