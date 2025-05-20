package app

import (
	"fmt"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

type DistributedSystemDemo struct {
	SenderNodeName   string
	ReceiverNodeName string
	NumberOfProcess  int

	senderNode   gen.Node
	receiverNode gen.Node
	cookie       string
}

func NewDistributedSystemDemo(SenderNodeName string, ReceiverNodeName string, NumberOfProcess int) *DistributedSystemDemo {
	return &DistributedSystemDemo{
		SenderNodeName:   SenderNodeName,
		ReceiverNodeName: ReceiverNodeName,
		NumberOfProcess:  NumberOfProcess,
		cookie:           "123",
	}
}

func (app *DistributedSystemDemo) initSenderNode() {
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

	for i := 1; i <= app.NumberOfProcess; i++ {
		app.senderNode.SpawnRegister(gen.Atom(fmt.Sprintf("sender_%d", i)), FactorySenderActor, gen.ProcessOptions{}, SenderActorParams{
			ReceiverName:     fmt.Sprintf("receiver_%d", i),
			ReceiverNodeName: app.ReceiverNodeName,
		})
		app.senderNode.SpawnRegister(gen.Atom(fmt.Sprintf("receiver_%d", i)), FactoryReceiverActor, gen.ProcessOptions{})
	}
}

func (app *DistributedSystemDemo) initReceiverNode() {
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

	for i := 1; i <= app.NumberOfProcess; i++ {
		app.receiverNode.SpawnRegister(gen.Atom(fmt.Sprintf("receiver_%d", i)), FactoryReceiverActor, gen.ProcessOptions{})
	}
}

func (app *DistributedSystemDemo) Start() {
	app.initSenderNode()
	app.initReceiverNode()

	fmt.Println()
	fmt.Println()

	app.senderNode.Wait()
}
