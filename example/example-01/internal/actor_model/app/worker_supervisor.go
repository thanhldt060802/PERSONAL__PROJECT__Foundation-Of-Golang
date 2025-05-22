package app

import (
	"fmt"
	"thanhldt060802/internal/actor_model/types"
	"thanhldt060802/internal/repository"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type WorkerSupervisor struct {
	act.Supervisor

	taskRepository         repository.TaskRepository
	numberOfInitialWorkers int

	availableWorkerMap map[string]bool
}

func FactoryWorkerSupervisor() gen.ProcessBehavior {
	return &WorkerSupervisor{}
}

func (workerSupervisor *WorkerSupervisor) Init(args ...any) (act.SupervisorSpec, error) {
	workerSupervisor.taskRepository = args[0].(repository.TaskRepository)
	workerSupervisor.numberOfInitialWorkers = args[1].(int)
	workerSupervisor.availableWorkerMap = map[string]bool{}

	supervisorSpec := act.SupervisorSpec{}
	supervisorSpec.EnableHandleChild = true
	supervisorSpec.DisableAutoShutdown = true
	supervisorSpec.Type = act.SupervisorTypeOneForOne
	supervisorSpec.Restart.Strategy = act.SupervisorStrategyTransient
	supervisorSpec.Restart.Intensity = 100
	supervisorSpec.Restart.Period = 5

	supervisorSpec.Children = []act.SupervisorChildSpec{}
	for i := 1; i <= workerSupervisor.numberOfInitialWorkers; i++ {
		supervisorSpec.Children = append(supervisorSpec.Children, act.SupervisorChildSpec{
			Name:    gen.Atom(fmt.Sprintf("worker_%d", i)),
			Factory: FactoryWorkerActor,
			Options: gen.ProcessOptions{},
			Args:    nil,
		})
		workerSupervisor.availableWorkerMap[fmt.Sprintf("worker_%d", i)] = false
	}

	workerSupervisor.Log().Info("Started worker supervisor %s %s on %s", workerSupervisor.PID(), workerSupervisor.Name(), workerSupervisor.Node().Name())
	return supervisorSpec, nil
}

func (workerSupervisor *WorkerSupervisor) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	workerSupervisor.Log().Info("Actor start with name %s and PID %s", childName, pid)
	return nil
}

func (workerSupervisor *WorkerSupervisor) HandleChildTerminate(name gen.Atom, pid gen.PID, reason error) error {
	if reason.Error() == gen.TerminateReasonNormal.Error() {
		workerName := name.String()
		workerName = workerName[1 : len(workerName)-1]
		workerSupervisor.availableWorkerMap[workerName] = true
	}

	return nil
}

func (workerSupervisor *WorkerSupervisor) HandleMessage(from gen.PID, message any) error {
	switch receivedMessage := message.(type) {
	case types.DoStart:
		{
			for _, supervisorChildSpec := range workerSupervisor.Children() {
				workerSupervisor.Send(supervisorChildSpec.Name, types.DoStart{})
			}

			return nil
		}
	case types.GetExistedWorkersMessage:
		{
			workerSupervisor.getExistedWorkers(receivedMessage)
			return nil
		}
	case types.RunTaskMessage:
		{
			workerSupervisor.runTask(receivedMessage)
			return nil
		}
	}

	return nil
}

func (workerSupervisor *WorkerSupervisor) getExistedWorkers(message types.GetExistedWorkersMessage) {
	workerNames := []string{}
	running := []string{}
	available := []string{}

	for _, supervisorChildSpec := range workerSupervisor.Children() {
		workerName := supervisorChildSpec.Name.String()
		workerName = workerName[1 : len(workerName)-1]
		workerNames = append(workerNames, workerName)
		if workerSupervisor.availableWorkerMap[workerName] {
			available = append(available, workerName)
		} else {
			running = append(running, workerName)
		}
	}
	message.WorkerNames <- workerNames
	message.Running <- running
	message.Available <- available
}

func (workerSupervisor *WorkerSupervisor) runTask(message types.RunTaskMessage) {
	if available, ok := workerSupervisor.availableWorkerMap[message.WorkerName]; ok {
		if available {
			workerSupervisor.Log().Info("--> Restart existed actor %s", message.WorkerName)
			if err := workerSupervisor.StartChild(gen.Atom(message.WorkerName), workerSupervisor.taskRepository, message.TaskId); err != nil {
				workerSupervisor.Log().Error("--- Restart existed actor %s failed: %s", message.WorkerName, err.Error())
			}
			workerSupervisor.availableWorkerMap[message.WorkerName] = false
			workerSupervisor.Log().Info("--> Restart existed actor %s successful", message.WorkerName)
		} else {
			workerSupervisor.Log().Warning("--- Actor %s is running", message.WorkerName)
		}
	} else {
		workerSupervisor.Log().Info("--> Start new actor %s", message.WorkerName)
		if err := workerSupervisor.AddChild(act.SupervisorChildSpec{
			Name:    gen.Atom(message.WorkerName),
			Factory: FactoryWorkerActor,
			Options: gen.ProcessOptions{},
			Args:    []any{workerSupervisor.taskRepository, message.TaskId},
		}); err != nil {
			workerSupervisor.Log().Info("--> Start new actor %s failed: %s", message.WorkerName, err.Error())
		}
		workerSupervisor.availableWorkerMap[message.WorkerName] = false
		workerSupervisor.Log().Info("--> Start new actor %s successful", message.WorkerName)
	}
}
