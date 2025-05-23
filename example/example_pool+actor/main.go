package main

import (
	"math/rand"
	"thanhldt060802/app"
	"time"

	"ergo.services/application/observer"
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

	nodeOptions.Applications = []gen.ApplicationBehavior{
		observer.CreateApp(observer.Options{}),
	}
	nodeOptions.Log.DefaultLogger.Disable = true
	nodeOptions.Log.Loggers = append(nodeOptions.Log.Loggers, gen.Logger{Name: "mynode", Logger: loggerColored})

	node, err := ergo.StartNode(gen.Atom("mynode@localhost"), nodeOptions)
	if err != nil {
		panic(err)
	}

	node.Spawn(app.FactoryWorkerPool, gen.ProcessOptions{})

	select {}

}
