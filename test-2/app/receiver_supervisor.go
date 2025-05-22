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
		receiverSupervisor.SpawnRegister(gen.Atom(fmt.Sprintf("receiver_%d", i)), FactoryReceiverActor, gen.ProcessOptions{}, types.Task{Task: fmt.Sprintf("task-%d", i), Situation: "done"})
	}

	return supervisorSpec, nil
}

func (receiverSupervisor *ReceiverSupervisor) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	receiverSupervisor.Log().Info("Actor start with name %s and PID %s", childName, pid)
	return nil
}

func (receiverSupervisor *ReceiverSupervisor) HandleChildTerminate(name gen.Atom, pid gen.PID, reason error) error {
	fmt.Println()
	for _, supervisorChildSpec := range receiverSupervisor.Children() {
		fmt.Printf("%s\n", supervisorChildSpec.Name)
	}
	fmt.Println()
	return nil
}

func (receiverSupervisor *ReceiverSupervisor) HandleMessage(from gen.PID, message any) error {
	receivedMessage := message.(types.TaskMessage)

	receiverSupervisor.Log().Info("Start actor %s", receivedMessage.Receiver)

	if err := receiverSupervisor.AddChild(act.SupervisorChildSpec{
		Name:    gen.Atom(receivedMessage.Receiver),
		Factory: FactoryReceiverActor,
		Options: gen.ProcessOptions{},
		Args: []any{
			types.Task{Task: receivedMessage.Task, Situation: receivedMessage.Situation},
		},
	}); err != nil {
		receiverSupervisor.Log().Info(err.Error())
		if err := receiverSupervisor.StartChild(gen.Atom(receivedMessage.Receiver), types.Task{Task: receivedMessage.Task, Situation: receivedMessage.Situation}); err != nil {
			receiverSupervisor.Log().Info(err.Error())
		}
	}

	return nil
}
