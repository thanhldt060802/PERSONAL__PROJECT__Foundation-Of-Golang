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
}

func FactoryWorkerSupervisor() gen.ProcessBehavior {
	return &WorkerSupervisor{}
}

func (workerSupervisor *WorkerSupervisor) Init(args ...any) (act.SupervisorSpec, error) {
	workerSupervisor.taskRepository = args[0].(repository.TaskRepository)
	workerSupervisor.numberOfInitialWorkers = args[1].(int)

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
	}

	workerSupervisor.Log().Info("Started worker supervisor %s %s on %s", workerSupervisor.PID(), workerSupervisor.Name(), workerSupervisor.Node().Name())
	return supervisorSpec, nil
}

func (workerSupervisor *WorkerSupervisor) HandleChildStart(childName gen.Atom, pid gen.PID) error {
	workerSupervisor.Log().Info("Actor start with name %s and PID %s", childName, pid)
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
	case types.RunTasksMessage:
		{
			workerSupervisor.runTasks(receivedMessage)
			return nil
		}
	}

	return nil
}

func (workerSupervisor *WorkerSupervisor) getExistedWorkers(message types.GetExistedWorkersMessage) {
	workerNames := []string{}
	for _, supervisorChildSpec := range workerSupervisor.Children() {
		workerName := supervisorChildSpec.Name.String()
		workerNames = append(workerNames, workerName[1:len(workerName)-1])
	}
	message.WorkerNames <- workerNames
}

func (workerSupervisor *WorkerSupervisor) runTask(message types.RunTaskMessage) error {
	if err := workerSupervisor.AddChild(act.SupervisorChildSpec{
		Name:    gen.Atom(message.WorkerName),
		Factory: FactoryWorkerActor,
		Options: gen.ProcessOptions{},
		Args:    []any{workerSupervisor.taskRepository, message.TaskId},
	}); err != nil {
		workerSupervisor.Log().Warning("--> Restart exsited actor %s", message.WorkerName)
		if err := workerSupervisor.StartChild(gen.Atom(message.WorkerName), workerSupervisor.taskRepository, message.TaskId); err != nil {
			workerSupervisor.Log().Error("--- Restart error: %s", err.Error())
			return err
		}
	} else {
		workerSupervisor.Log().Info("--> Start new actor %s", message.WorkerName)
	}

	return nil
}

func (workerSupervisor *WorkerSupervisor) runTasks(message types.RunTasksMessage) {
	numberOfTasks := len(message.TaskIds)
	numberOfExistedWorkers := len(workerSupervisor.Children())
	index := 0

	for _, supervisorChildSpec := range workerSupervisor.Children() {
		if index >= numberOfTasks {
			break
		}
		workerName := supervisorChildSpec.Name.String()
		workerName = workerName[1 : len(workerName)-1]
		dispatchingMessage := types.RunTaskMessage{
			WorkerName: workerName,
			TaskId:     message.TaskIds[index],
		}
		if err := workerSupervisor.runTask(dispatchingMessage); err == nil {
			index++
		}
	}

	for index < numberOfTasks {
		numberOfExistedWorkers++
		workerName := fmt.Sprintf("worker_%d", numberOfExistedWorkers)
		dispatchingMessage := types.RunTaskMessage{
			WorkerName: workerName,
			TaskId:     message.TaskIds[index],
		}
		workerSupervisor.runTask(dispatchingMessage)
		index++
	}
}
