package app

import (
	"fmt"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

type DistributedSystemHandleTask struct {
	DispatcherNodeName string
	ReceiverNodeName   string
	NumberOfProcess    int

	dispatcherNode gen.Node
	receiverNode   gen.Node
	cookie         string
}

func NewDistributedSystemHandleTask(DispatcherNodeName string, ReceiverNodeName string, NumberOfProcess int) *DistributedSystemHandleTask {
	return &DistributedSystemHandleTask{
		DispatcherNodeName: DispatcherNodeName,
		ReceiverNodeName:   ReceiverNodeName,
		NumberOfProcess:    NumberOfProcess,
		cookie:             "123",
	}
}

func (app *DistributedSystemHandleTask) initDispatcherNode() {
	var nodeOptions gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	nodeOptions.Log.DefaultLogger.Disable = true
	nodeOptions.Log.Loggers = append(nodeOptions.Log.Loggers, gen.Logger{Name: app.DispatcherNodeName, Logger: loggerColored})
	nodeOptions.Network.Cookie = app.cookie

	node, err := ergo.StartNode(gen.Atom(app.DispatcherNodeName), nodeOptions)
	if err != nil {
		panic(err)
	}
	app.dispatcherNode = node

	app.dispatcherNode.SpawnRegister("dispatcher", FactoryDispatcherActor, gen.ProcessOptions{})
}

func (app *DistributedSystemHandleTask) initReceiverNode() {
	var nodeOptions gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	nodeOptions.Log.DefaultLogger.Disable = true
	nodeOptions.Log.Loggers = append(nodeOptions.Log.Loggers, gen.Logger{Name: app.ReceiverNodeName, Logger: loggerColored})
	nodeOptions.Network.Cookie = app.cookie

	node, err := ergo.StartNode(gen.Atom(app.ReceiverNodeName), nodeOptions)
	if err != nil {
		panic(err)
	}
	app.receiverNode = node

	receiverSupervisorPID, _ := app.receiverNode.Spawn(FactoryReceiverSupervisor, gen.ProcessOptions{}, ReceiverSupervisorParams{
		dispatcherProcessName: "dispatcher",
		dispatcherNodeName:    app.DispatcherNodeName,
		numberOfProcess:       app.NumberOfProcess,
	})
	app.receiverNode.Log().Info("Supervisor for receiver node is started succesfully with PID %s", receiverSupervisorPID)
}

func (app *DistributedSystemHandleTask) Start() {
	app.initDispatcherNode()
	app.initReceiverNode()

	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	app.receiverNode.Wait()
}
