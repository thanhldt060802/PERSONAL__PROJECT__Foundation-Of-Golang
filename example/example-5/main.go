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

	var nodeOptions gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	nodeOptions.Log.DefaultLogger.Disable = true
	nodeOptions.Log.Loggers = append(nodeOptions.Log.Loggers, gen.Logger{Name: "my_node", Logger: loggerColored})
	nodeOptions.Network.Cookie = "123"

	myNode, err := ergo.StartNode(gen.Atom("mynode@localhost"), nodeOptions)
	if err != nil {
		panic(err)
	}

	myNode.SpawnRegister(gen.Atom("receiver_fsm"), app.FactoryReceiverFSMActor, gen.ProcessOptions{})

	fmt.Println()
	fmt.Println()

	myNode.Wait()
}
