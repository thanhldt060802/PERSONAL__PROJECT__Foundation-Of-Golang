package app

import (
	"fmt"
	"thanhtldt060802/actor_model/types"
	"thanhtldt060802/internal/repository"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverSupervisorParams struct {
	taskRepository         repository.TaskRepository
	numberOfInitialProcess int
}

type ReceiverSupervisor struct {
	act.Supervisor

	params ReceiverSupervisorParams
}

func FactoryReceiverSupervisor() gen.ProcessBehavior {
	return &ReceiverSupervisor{}
}

func (receiverSupervisor *ReceiverSupervisor) Init(args ...any) (act.SupervisorSpec, error) {
	receiverSupervisor.params = args[0].(ReceiverSupervisorParams)

	supervisorSpec := act.SupervisorSpec{}
	supervisorSpec.EnableHandleChild = true
	supervisorSpec.DisableAutoShutdown = true
	supervisorSpec.Type = act.SupervisorTypeOneForOne
	supervisorSpec.Children = []act.SupervisorChildSpec{}
	supervisorSpec.Restart.Strategy = act.SupervisorStrategyTransient
	supervisorSpec.Restart.Intensity = 10
	supervisorSpec.Restart.Period = 5

	for i := 1; i <= receiverSupervisor.params.numberOfInitialProcess; i++ {
		supervisorSpec.Children = append(supervisorSpec.Children, act.SupervisorChildSpec{
			Name:    gen.Atom(fmt.Sprintf("receiver_%d", i)),
			Factory: FactoryReceiverActor,
			Options: gen.ProcessOptions{},
			Args:    nil,
		})
	}

	return supervisorSpec, nil
}

func (receiverSupervisor *ReceiverSupervisor) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	receiverSupervisor.Log().Info("Actor start with name %s and PID %s", childName, pid)
	return nil
}

func (receiverSupervisor *ReceiverSupervisor) HandleMessage(from gen.PID, message any) error {
	if message == nil {
		for _, supervisorChildSpec := range receiverSupervisor.Children() {
			receiverSupervisor.Send(supervisorChildSpec.Name, nil)
		}
	} else {
		switch receivedMessage := message.(type) {
		case types.ExistedReceiverNamesMessage:
			{
				existedReceiverNames := []string{}
				for _, supervisorChildSpec := range receiverSupervisor.Children() {
					existedReceiverNames = append(existedReceiverNames, supervisorChildSpec.Name.String())
				}
				receivedMessage.ReceiverNames <- existedReceiverNames

				return nil
			}
		case types.NewTaskMessage:
			{
				newTask := types.Task{
					TaskId: receivedMessage.TaskId,
				}

				if err := receiverSupervisor.AddChild(act.SupervisorChildSpec{
					Name:    gen.Atom(receivedMessage.Receiver),
					Factory: FactoryReceiverActor,
					Options: gen.ProcessOptions{},
					Args: []any{
						ReceiverActorParams{task: newTask},
					},
				}); err != nil {
					receiverSupervisor.Log().Warning("Restart exsited actor %s", receivedMessage.Receiver)
					if err := receiverSupervisor.StartChild(gen.Atom(receivedMessage.Receiver), ReceiverActorParams{taskRepository: receiverSupervisor.params.taskRepository, task: newTask}); err != nil {
						receiverSupervisor.Log().Error("Restart error -- %s", err.Error())
					}
				} else {
					receiverSupervisor.Log().Info("Start new actor %s", receivedMessage.Receiver)
				}

				return nil
			}
		}
	}

	return nil
}
