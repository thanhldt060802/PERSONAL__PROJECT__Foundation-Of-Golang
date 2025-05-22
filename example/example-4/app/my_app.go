package app

import (
	"fmt"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

type DistributedSystemHandleTask struct {
	SenderNodeName   string
	ReceiverNodeName string
	NumberOfProcess  int

	senderNode   gen.Node
	receiverNode gen.Node
	cookie       string
}

func NewDistributedSystemHandleTask(SenderNodeName string, ReceiverNodeName string, NumberOfProcess int) *DistributedSystemHandleTask {
	return &DistributedSystemHandleTask{
		SenderNodeName:   SenderNodeName,
		ReceiverNodeName: ReceiverNodeName,
		NumberOfProcess:  NumberOfProcess,
		cookie:           "123",
	}
}

func (app *DistributedSystemHandleTask) initSenderNode() {
	var nodeOptions gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	nodeOptions.Log.DefaultLogger.Disable = true
	nodeOptions.Log.Loggers = append(nodeOptions.Log.Loggers, gen.Logger{Name: app.SenderNodeName, Logger: loggerColored})
	nodeOptions.Network.Cookie = app.cookie

	node, err := ergo.StartNode(gen.Atom(app.SenderNodeName), nodeOptions)
	if err != nil {
		panic(err)
	}
	app.senderNode = node

	app.senderNode.SpawnRegister("sender", FactorySenderActor, gen.ProcessOptions{})
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

	receiverSupervisorPID, _ := app.receiverNode.Spawn(FactoryReceiverSupervisor, gen.ProcessOptions{}, ReceiverSupervisorParam{
		SenderName:      "sender",
		SenderNodeName:  app.SenderNodeName,
		NumberOfProcess: app.NumberOfProcess,
	})
	app.receiverNode.Log().Info("Supervisor for receiver node is started succesfully with PID %s", receiverSupervisorPID)
}

func (app *DistributedSystemHandleTask) Start() {
	app.initSenderNode()
	app.initReceiverNode()

	fmt.Println()
	fmt.Println()

	app.receiverNode.Wait()
}
