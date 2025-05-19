package app

import (
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
	var supervisorSpec act.SupervisorSpec
	supervisorSpec.EnableHandleChild = true
	supervisorSpec.Type = act.SupervisorTypeOneForOne
	supervisorSpec.Children = []act.SupervisorChildSpec{
		{
			Name:    gen.Atom("actor1"),
			Factory: FactoryReceiverFSMActor,
			Options: gen.ProcessOptions{},
		},
		{
			Name:    gen.Atom("actor2"),
			Factory: FactoryReceiverFSMActor,
			Options: gen.ProcessOptions{},
		},
	}
	supervisorSpec.Restart.Strategy = act.SupervisorStrategyTransient
	supervisorSpec.Restart.Intensity = 15
	supervisorSpec.Restart.Period = 5

	return supervisorSpec, nil
}

func (receiverSupervisor *ReceiverSupervisor) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	receiverSupervisor.Node().RegisterName(gen.Atom(childName), pid)
	return nil
}
