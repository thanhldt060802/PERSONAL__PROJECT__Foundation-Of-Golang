package app

import (
	"fmt"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverSupervisorParams struct {
	numberOfProcess int
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

	var supervisorSpec act.SupervisorSpec
	supervisorSpec.EnableHandleChild = true
	supervisorSpec.Type = act.SupervisorTypeOneForOne
	supervisorSpec.Children = []act.SupervisorChildSpec{}
	supervisorSpec.Restart.Strategy = act.SupervisorStrategyTransient
	supervisorSpec.Restart.Intensity = 15
	supervisorSpec.Restart.Period = 5

	for i := 1; i <= receiverSupervisor.params.numberOfProcess; i++ {
		supervisorSpec.Children = append(supervisorSpec.Children, act.SupervisorChildSpec{
			Name:    gen.Atom(fmt.Sprintf("receiver_fsm_%d", i)),
			Factory: FactoryReceiverFSMActor,
			Options: gen.ProcessOptions{},
		})
	}

	return supervisorSpec, nil
}

func (receiverSupervisor *ReceiverSupervisor) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	receiverSupervisor.Node().RegisterName(gen.Atom(childName), pid)
	return nil
}
