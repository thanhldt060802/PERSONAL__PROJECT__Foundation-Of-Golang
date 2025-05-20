package app

import (
	"fmt"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ReceiverSupervisorParam struct {
	SenderName      string
	SenderNodeName  string
	NumberOfProcess int
}

type ReceiverSupervisor struct {
	act.Supervisor

	params ReceiverSupervisorParam
}

func FactoryReceiverSupervisor() gen.ProcessBehavior {
	return &ReceiverSupervisor{}
}

func (receiverSupervisor *ReceiverSupervisor) Init(args ...any) (act.SupervisorSpec, error) {
	receiverSupervisor.params = args[0].(ReceiverSupervisorParam)

	var supervisorSpec act.SupervisorSpec
	supervisorSpec.EnableHandleChild = true
	supervisorSpec.Type = act.SupervisorTypeOneForOne
	supervisorSpec.Children = []act.SupervisorChildSpec{}
	supervisorSpec.Restart.Strategy = act.SupervisorStrategyTransient
	supervisorSpec.Restart.Intensity = 10
	supervisorSpec.Restart.Period = 5

	for i := 1; i <= receiverSupervisor.params.NumberOfProcess; i++ {
		supervisorSpec.Children = append(supervisorSpec.Children, act.SupervisorChildSpec{
			Name:    gen.Atom(fmt.Sprintf("receiver_%d", i)),
			Factory: FactoryReceiverActor,
			Options: gen.ProcessOptions{},
			Args: []any{
				ReceiverActorParam{
					SenderName:     receiverSupervisor.params.SenderName,
					SenderNodeName: receiverSupervisor.params.SenderNodeName,
				},
			},
		})
	}

	return supervisorSpec, nil
}

func (receiverSupervisor *ReceiverSupervisor) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	receiverSupervisor.Node().RegisterName(gen.Atom(childName), pid)
	return nil
}
