package main

import (
	"thanhldt060802/app"
	"time"

	"ergo.services/application/observer"
	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

func main() {

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
	nodeOptions.Log.Loggers = append(nodeOptions.Log.Loggers, gen.Logger{Name: "mynode", Logger: loggerColored})

	node, err := ergo.StartNode(gen.Atom("mynode@localhost"), nodeOptions)
	if err != nil {
		panic(err)
	}

	node.SpawnRegister("worker_1", app.FactoryWorkerActor, gen.ProcessOptions{})
	node.SpawnRegister("worker_2", app.FactoryWorkerActor, gen.ProcessOptions{})
	node.SpawnRegister("worker_3", app.FactoryWorkerActor, gen.ProcessOptions{})
	node.SpawnRegister("my_tcp", app.FactoryMyTCP, gen.ProcessOptions{})

	select {}

}
