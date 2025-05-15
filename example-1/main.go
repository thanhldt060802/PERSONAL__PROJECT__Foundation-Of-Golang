package main

import (
	"math/rand"
	"thanhldt060802/actors"
	"thanhldt060802/common"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var node1Options, node2Options gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	node1Options.Log.DefaultLogger.Disable = true
	node1Options.Log.Loggers = append(node1Options.Log.Loggers, gen.Logger{Name: "node-1", Logger: loggerColored})
	node1Options.Network.Cookie = "123"

	node2Options.Log.DefaultLogger.Disable = true
	node2Options.Log.Loggers = append(node2Options.Log.Loggers, gen.Logger{Name: "node-2", Logger: loggerColored})
	node2Options.Network.Cookie = "123"

	node1, err := ergo.StartNode("node1@localhost", node1Options)
	if err != nil {
		panic(err)
	}

	node1.SpawnRegister("sender-1", actors.FactorySenderActor, gen.ProcessOptions{})
	node1.SpawnRegister("sender-2", actors.FactorySenderActor, gen.ProcessOptions{})
	node1.SpawnRegister("receiver-1", actors.FactoryReceiverActor, gen.ProcessOptions{})

	node2, err := ergo.StartNode("node2@localhost", node2Options)
	if err != nil {
		panic(err)
	}

	node2.SpawnRegister("receiver-1", actors.FactoryReceiverActor, gen.ProcessOptions{})

	time.Sleep(2 * time.Second)

	node1.Send(gen.Atom("sender-1"), common.DoCallLocal{})
	node1.Send(gen.Atom("sender-2"), common.DoCallLocal{})

	node1.Wait()
}
