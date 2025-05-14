package main

import (
	"math/rand"
	"thanhldt060802/actors"
	"thanhldt060802/common"
	"thanhldt060802/supervisors"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var node1Options, node2Options gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true, IncludeBehavior: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	node1Options.Log.DefaultLogger.Disable = true
	node1Options.Log.Loggers = append(node1Options.Log.Loggers, gen.Logger{Name: "node 1", Logger: loggerColored})
	node1Options.Network.Cookie = "123"

	node2Options.Log.DefaultLogger.Disable = true
	node2Options.Log.Loggers = append(node2Options.Log.Loggers, gen.Logger{Name: "node 2", Logger: loggerColored})
	node2Options.Network.Cookie = "123"

	node1, err := ergo.StartNode("node1@localhost", node1Options)
	if err != nil {
		panic(err)
	}

	node1.SpawnRegister("a", actors.FactoryActorA, gen.ProcessOptions{})
	node1.SpawnRegister("b", actors.FactoryActorB, gen.ProcessOptions{})
	node1.SpawnRegister("c1", actors.FactoryActorC, gen.ProcessOptions{})

	node2, err := ergo.StartNode("node2@localhost", node2Options)
	if err != nil {
		panic(err)
	}

	// node2.SpawnRegister("c1", actors.FactoryActorC, gen.ProcessOptions{})
	supervisorCPID, _ := node2.Spawn(supervisors.FactorySupervisorSpecC, gen.ProcessOptions{})
	node2.Log().Info("Supervisor for node 2 is started succesfully with pid %s", supervisorCPID)

	time.Sleep(2 * time.Second)

	node1.Send(gen.Atom("a"), common.DoCallLocal{})
	node1.Send(gen.Atom("b"), common.DoCallLocal{})

	node1.Wait()
}
