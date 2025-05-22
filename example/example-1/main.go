package main

import (
	"fmt"
	"math/rand"
	"thanhldt060802/app"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Set options for node1 and node2

	var node1Options, node2Options gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	node1Options.Log.DefaultLogger.Disable = true
	node1Options.Log.Loggers = append(node1Options.Log.Loggers, gen.Logger{Name: "node1", Logger: loggerColored})
	node1Options.Network.Cookie = "123"

	node2Options.Log.DefaultLogger.Disable = true
	node2Options.Log.Loggers = append(node2Options.Log.Loggers, gen.Logger{Name: "node2", Logger: loggerColored})
	node2Options.Network.Cookie = "123"

	// Init node1 which owns sender_1, sender_2 and local receiver_1

	node1, err := ergo.StartNode(gen.Atom("node1@localhost"), node1Options)
	if err != nil {
		panic(err)
	}

	node1.SpawnRegister(gen.Atom("sender_1"), app.FactorySenderActor, gen.ProcessOptions{}, app.SenderActorParams{
		ReceiverName:     "receiver_1",
		ReceiverNodeName: "node2@localhost",
	})
	node1.SpawnRegister(gen.Atom("sender_2"), app.FactorySenderActor, gen.ProcessOptions{}, app.SenderActorParams{
		ReceiverName:     "receiver_1",
		ReceiverNodeName: "node2@localhost",
	})
	node1.SpawnRegister(gen.Atom("receiver_1"), app.FactoryReceiverActor, gen.ProcessOptions{})

	// Init node2 which owns receiver_1

	node2, err := ergo.StartNode(gen.Atom("node2@localhost"), node2Options)
	if err != nil {
		panic(err)
	}

	node2.SpawnRegister(gen.Atom("receiver_1"), app.FactoryReceiverActor, gen.ProcessOptions{})

	// Trigger sender_1 and sender_2 on node1

	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	node1.Send(gen.Atom("sender_1"), "local")
	node1.Send(gen.Atom("sender_2"), "local")

	node1.Wait()
}
