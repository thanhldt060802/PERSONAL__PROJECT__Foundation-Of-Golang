package app

import (
	"fmt"
	"thanhtldt060802/internal/repository"
	"time"

	"ergo.services/application/observer"
	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

type MyApp struct {
	taskRepository         repository.TaskRepository
	nodeName               string
	numberOfInitialProcess int

	node          gen.Node
	supervisorPID gen.PID
}

func New(taskRepository repository.TaskRepository, nodeName string, numberOfInitialProcess int) *MyApp {
	return &MyApp{
		taskRepository:         taskRepository,
		nodeName:               nodeName,
		numberOfInitialProcess: numberOfInitialProcess,
	}
}

func (myApp *MyApp) Start() {
	var nodeOptions gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	nodeOptions.Applications = []gen.ApplicationBehavior{
		observer.CreateApp(observer.Options{}),
	}
	nodeOptions.Log.DefaultLogger.Disable = true
	nodeOptions.Log.Loggers = append(nodeOptions.Log.Loggers, gen.Logger{Name: myApp.nodeName, Logger: loggerColored})

	node, err := ergo.StartNode(gen.Atom(fmt.Sprintf("%s@localhost", myApp.nodeName)), nodeOptions)
	if err != nil {
		panic(err)
	}

	myApp.node = node

	supervisorPID, _ := myApp.node.Spawn(FactoryReceiverSupervisor, gen.ProcessOptions{}, ReceiverSupervisorParams{taskRepository: myApp.taskRepository, numberOfInitialProcess: myApp.numberOfInitialProcess})

	myApp.supervisorPID = supervisorPID

	myApp.node.Send(supervisorPID, nil)
}

func (myApp *MyApp) Node() gen.Node {
	return myApp.node
}

func (myApp *MyApp) SupervisorPID() gen.PID {
	return myApp.supervisorPID
}
