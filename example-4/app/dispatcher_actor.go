package app

import (
	"context"
	"fmt"
	"thanhldt060802/repository"
	"thanhldt060802/types"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type DispatcherActor struct {
	act.Actor
}

func FactoryDispatcherActor() gen.ProcessBehavior {
	return &DispatcherActor{}
}

func (dispatcherActor *DispatcherActor) Init(args ...any) error {
	dispatcherActor.Log().Info("started process %s %s on %s", dispatcherActor.PID(), dispatcherActor.Name(), dispatcherActor.Node().Name())
	return nil
}

func (dispatcherActor *DispatcherActor) HandleMessage(from gen.PID, message any) error {
	receivedMessage := message.(types.ReturnTaskMessage)

	switch receivedMessage.SelfStatus {
	case "idle":
		{
			task, err := repository.TaskRepositoryInstance.GetAvailable(context.Background())
			if err != nil {
				return nil
			}

			sendingMessage := types.DoTaskMessage{
				Id:       task.Id,
				Progress: task.Progress,
				Target:   task.Target,
			}

			dispatcherActor.Log().Info("--> %s: %#v", from, sendingMessage)
			dispatcherActor.Send(from, sendingMessage)

			return nil
		}
	case "complete":
		{
			fmt.Println("COMPLETE")
			foundTask, err := repository.TaskRepositoryInstance.GetById(context.Background(), receivedMessage.Id)
			if err != nil {
				return fmt.Errorf("processing failed: %s", err.Error())
			}
			foundTask.Progress = receivedMessage.Progress
			foundTask.Status = "COMPLETE"

			if err := repository.TaskRepositoryInstance.Update(context.Background(), foundTask); err != nil {
				return fmt.Errorf("processing failed: %s", err.Error())
			}
			dispatcherActor.Log().Info("--- Update task successful by %s'returning: %#v", from, receivedMessage)

			return nil
		}
		// case "cancel":
		// 	{
		// 		foundTask, err := repository.TaskRepositoryInstance.GetById(context.Background(), receivedMessage.Id)
		// 		if err != nil {
		// 			return fmt.Errorf("processing failed: %s", err.Error())
		// 		}
		// 		foundTask.Progress = receivedMessage.Progress
		// 		foundTask.Status = "CANCEL"

		// 		if err := repository.TaskRepositoryInstance.Update(context.Background(), foundTask); err != nil {
		// 			return fmt.Errorf("processing failed: %s", err.Error())
		// 		}
		// 		dispatcherActor.Log().Info("--- Update task successful by %s'returning: %#v", from, receivedMessage)

		// 		return nil
		// 	}
	}

	dispatcherActor.Log().Error("--- Unknown message %#v", message)
	return nil
}

func (dispatcherActor *DispatcherActor) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	receivedMessage := request.(types.ReturnTaskMessage)

	switch receivedMessage.SelfStatus {
	case "cancel":
		{
			foundTask, err := repository.TaskRepositoryInstance.GetById(context.Background(), receivedMessage.Id)
			if err != nil {
				return nil, fmt.Errorf("processing failed: %s", err.Error())
			}
			foundTask.Progress = receivedMessage.Progress
			foundTask.Status = "CANCEL"

			if err := repository.TaskRepositoryInstance.Update(context.Background(), foundTask); err != nil {
				return nil, fmt.Errorf("processing failed: %s", err.Error())
			}
			dispatcherActor.Log().Info("--- Update task successful by %s'returning: %#v", from, receivedMessage)

			return nil, nil
		}
	}

	dispatcherActor.Log().Error("--- Unknown request %#v", request)
	return nil, nil
}
