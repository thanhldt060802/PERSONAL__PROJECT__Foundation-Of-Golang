package app

import (
	"fmt"
	"thanhtldt060802/types"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverSupervisor struct {
	act.Supervisor
}

func FactoryReceiverSupervisor() gen.ProcessBehavior {
	return &ReceiverSupervisor{}
}

func (receiverSupervisor *ReceiverSupervisor) Init(args ...any) (act.SupervisorSpec, error) {
	supervisorSpec := act.SupervisorSpec{}
	supervisorSpec.EnableHandleChild = true
	supervisorSpec.DisableAutoShutdown = true
	supervisorSpec.Type = act.SupervisorTypeOneForOne
	supervisorSpec.Children = []act.SupervisorChildSpec{}
	supervisorSpec.Restart.Strategy = act.SupervisorStrategyTransient
	supervisorSpec.Restart.Intensity = 10
	supervisorSpec.Restart.Period = 5

	for i := 1; i <= 5; i++ {
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

func (receiverSupervisor *ReceiverSupervisor) HandleChildTerminate(name gen.Atom, pid gen.PID, reason error) error {
	fmt.Println(name)
	fmt.Println()
	for _, supervisorChildSpec := range receiverSupervisor.Children() {
		fmt.Printf("%s\n", supervisorChildSpec.Name)
	}
	fmt.Println()
	return nil
}

func (receiverSupervisor *ReceiverSupervisor) HandleMessage(from gen.PID, message any) error {
	if message == nil {
		for _, supervisorChildSpec := range receiverSupervisor.Children() {
			receiverSupervisor.Send(supervisorChildSpec.Name, nil)
		}
	} else {
		receivedMessage := message.(types.TaskMessage)

		if err := receiverSupervisor.AddChild(act.SupervisorChildSpec{
			Name:    gen.Atom(receivedMessage.Receiver),
			Factory: FactoryReceiverActor,
			Options: gen.ProcessOptions{},
			Args: []any{
				types.Task{Task: receivedMessage.Task, Situation: receivedMessage.Situation},
			},
		}); err != nil {
			receiverSupervisor.Log().Warning("Restart exsited actor %s", receivedMessage.Receiver)
			if err := receiverSupervisor.StartChild(gen.Atom(receivedMessage.Receiver), types.Task{Task: receivedMessage.Task, Situation: receivedMessage.Situation}); err != nil {
				receiverSupervisor.Log().Error("Restart error -- %s", err.Error())
			}
		} else {
			receiverSupervisor.Log().Info("Start new actor %s", receivedMessage.Receiver)
		}
	}

	return nil
}
