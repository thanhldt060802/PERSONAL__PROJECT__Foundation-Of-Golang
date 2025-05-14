package supervisors

import (
	"thanhldt060802/actors"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type SupervisorSpecC struct {
	act.Supervisor
}

func FactorySupervisorSpecC() gen.ProcessBehavior {
	return &SupervisorSpecC{}
}

func (supervisorSpecC *SupervisorSpecC) Init(args ...any) (act.SupervisorSpec, error) {
	var supervisorSpec act.SupervisorSpec
	supervisorSpec.EnableHandleChild = true
	supervisorSpec.Type = act.SupervisorTypeOneForOne
	supervisorSpec.Children = []act.SupervisorChildSpec{
		{
			Name:    "c1",
			Factory: actors.FactoryActorC,
			Options: gen.ProcessOptions{},
		},
		// {
		// 	Name:    "c2",
		// 	Factory: actors.FactoryActorC,
		// 	Options: gen.ProcessOptions{},
		// },
		// {
		// 	Name:    "c3",
		// 	Factory: actors.FactoryActorC,
		// 	Options: gen.ProcessOptions{},
		// },
	}
	supervisorSpec.Restart.Strategy = act.SupervisorStrategyTransient
	supervisorSpec.Restart.Intensity = 3
	supervisorSpec.Restart.Period = 3

	return supervisorSpec, nil
}

func (supervisorSpecC *SupervisorSpecC) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	supervisorSpecC.Node().RegisterName(gen.Atom(childName), pid)
	supervisorSpecC.Log().Info("Registered process %s with pid %s", childName, pid)
	return nil
}
