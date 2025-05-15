package receiversupervisors

import (
	"thanhldt060802/actors"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverSupervisor struct {
	act.Supervisor
}

func FactoryReceiverSupervisor() gen.ProcessBehavior {
	return &ReceiverSupervisor{}
}

func (supervisorSpecC *ReceiverSupervisor) Init(args ...any) (act.SupervisorSpec, error) {
	var supervisorSpec act.SupervisorSpec
	supervisorSpec.EnableHandleChild = true
	supervisorSpec.Type = act.SupervisorTypeOneForOne
	supervisorSpec.Children = []act.SupervisorChildSpec{
		{
			Name:    "receiver_1",
			Factory: actors.FactoryReceiverActor,
			Options: gen.ProcessOptions{},
		},
		{
			Name:    "receiver_2",
			Factory: actors.FactoryReceiverActor,
			Options: gen.ProcessOptions{},
		},
		{
			Name:    "receiver_3",
			Factory: actors.FactoryReceiverActor,
			Options: gen.ProcessOptions{},
		},
		{
			Name:    "receiver_4",
			Factory: actors.FactoryReceiverActor,
			Options: gen.ProcessOptions{},
		},
		{
			Name:    "receiver_5",
			Factory: actors.FactoryReceiverActor,
			Options: gen.ProcessOptions{},
		},
	}
	supervisorSpec.Restart.Strategy = act.SupervisorStrategyTransient
	supervisorSpec.Restart.Intensity = 10
	supervisorSpec.Restart.Period = 10

	return supervisorSpec, nil
}

func (supervisorSpecC *ReceiverSupervisor) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	supervisorSpecC.Node().RegisterName(gen.Atom(childName), pid)
	return nil
}
