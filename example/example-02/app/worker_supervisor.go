package app

import (
	"fmt"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type WorkerSupervisor struct {
	act.Supervisor
}

func FactoryWorkerSupervisor() gen.ProcessBehavior {
	return &WorkerSupervisor{}
}

func (workerSupervisor *WorkerSupervisor) Init(args ...any) (act.SupervisorSpec, error) {
	supervisorSpec := act.SupervisorSpec{}
	supervisorSpec.EnableHandleChild = true
	supervisorSpec.DisableAutoShutdown = true
	supervisorSpec.Type = act.SupervisorTypeOneForOne
	supervisorSpec.Restart.Strategy = act.SupervisorStrategyTransient
	supervisorSpec.Restart.Intensity = 100
	supervisorSpec.Restart.Period = 5

	supervisorSpec.Children = []act.SupervisorChildSpec{}
	for i := 1; i <= 10; i++ {
		supervisorSpec.Children = append(supervisorSpec.Children, act.SupervisorChildSpec{
			Name:    gen.Atom(fmt.Sprintf("worker_%d", i)),
			Factory: FactoryWorkerActor,
			Options: gen.ProcessOptions{},
			Args:    []any{int64(i)},
		})
	}

	workerSupervisor.Log().Info("Started worker supervisor %s %s on %s", workerSupervisor.PID(), workerSupervisor.Name(), workerSupervisor.Node().Name())
	return supervisorSpec, nil
}

func (workerSupervisor *WorkerSupervisor) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	workerSupervisor.Log().Info("Actor start with name %s and PID %s", childName, pid)
	return nil
}

func (workerSupervisor *WorkerSupervisor) HandleChildTerminate(name gen.Atom, pid gen.PID, reason error) error {
	workerSupervisor.Log().Error("Actor %s terminated. Panic reason: %s", name, reason.Error())
	return nil
}
