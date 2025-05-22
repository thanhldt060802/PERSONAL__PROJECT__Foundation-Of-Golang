package app

import (
	"context"
	"math/rand"
	"thanhldt060802/repository"
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

	taskId       int64
	taskProgress int64
	taskTarget   int64
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
	receivedMessage := message.(types.DoTaskMessage)

	receiverActor.Log().Info("<-- %s: %#v", from, receivedMessage)

	receiverActor.taskId = receivedMessage.TaskId
	receiverActor.taskProgress = receivedMessage.TaskProgress
	receiverActor.taskTarget = receivedMessage.TaskTarget

	for receiverActor.taskProgress < receiverActor.taskTarget {
		if rand.Intn(10) == 0 {
			panic("Simulate crash")
		}

		time.Sleep(1 * time.Second)
		receiverActor.taskProgress++
	}

	foundTask, err := repository.TaskRepositoryInstance.GetById(context.Background(), receiverActor.taskId)
	foundTask.Progress = receiverActor.taskProgress
	foundTask.Status = "COMPLETE"
	if repository.TaskRepositoryInstance.Update(context.Background(), foundTask); err != nil {
		receiverActor.Log().Info("Update task failed: %#v (%s)", foundTask, err.Error())
	}
	receiverActor.Log().Info("Update task success: %#v", foundTask)

	receiverActor.Log().Info("Complete task %#v", foundTask)
	receiverActor.Send(from, "idle")

	return nil
}

func (receiverActor *ReceiverActor) Terminate(reason error) {
	receiverActor.Log().Error("Actor terminated. Panic reason: %s", reason.Error())

	foundTask, err := repository.TaskRepositoryInstance.GetById(context.Background(), receiverActor.taskId)
	foundTask.Progress = receiverActor.taskProgress
	foundTask.Status = "CANCEL"
	if repository.TaskRepositoryInstance.Update(context.Background(), foundTask); err != nil {
		receiverActor.Log().Info("Update task failed: %#v (%s)", foundTask, err.Error())
	}
	receiverActor.Log().Info("Update task success: %#v", foundTask)
}
