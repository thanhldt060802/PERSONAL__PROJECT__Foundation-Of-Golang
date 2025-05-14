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

func temp() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var nodeAOptions, nodeBOptions, nodeCOptions gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true, IncludeBehavior: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	// Setup for factory ActorA
	nodeAOptions.Log.DefaultLogger.Disable = true
	nodeAOptions.Log.Loggers = append(nodeAOptions.Log.Loggers, gen.Logger{Name: "node A", Logger: loggerColored})
	nodeAOptions.Network.Cookie = "123"

	// Setup for factory ActorB
	nodeBOptions.Log.DefaultLogger.Disable = true
	nodeBOptions.Log.Loggers = append(nodeBOptions.Log.Loggers, gen.Logger{Name: "node B", Logger: loggerColored})
	nodeBOptions.Network.Cookie = "123"

	// Setup for factory ActorC
	nodeCOptions.Log.DefaultLogger.TimeFormat = time.DateTime
	nodeCOptions.Log.DefaultLogger.IncludeName = true
	nodeCOptions.Log.DefaultLogger.IncludeBehavior = true
	nodeCOptions.Network.Cookie = "123"

	// Starting node A
	nodeA, err := ergo.StartNode("nodeA@localhost", nodeAOptions)
	if err != nil {
		panic(err)
	}

	// Starting node B
	nodeB, err := ergo.StartNode("nodeB@localhost", nodeBOptions)
	if err != nil {
		panic(err)
	}

	// Starting node C
	nodeC, err := ergo.StartNode("nodeC@localhost", nodeCOptions)
	if err != nil {
		panic(err)
	}
	defer nodeC.StopForce()

	// Handler

	// Spawn process 'c' is actor C on node C for remote request from other actor to 'c'
	// nodeC.SpawnRegister("c1", actors.FactoryActorC, gen.ProcessOptions{})
	// nodeC.SpawnRegister("c2", actors.FactoryActorC, gen.ProcessOptions{})
	supervisorCPID, _ := nodeC.Spawn(supervisors.FactorySupervisorSpecC, gen.ProcessOptions{})
	nodeC.Log().Info("Supervisor for node C is started succesfully with pid %s", supervisorCPID)

	// Caller

	// Spawm process 'a' is actor A on node A for triggering
	nodeA.SpawnRegister("a", actors.FactoryActorA, gen.ProcessOptions{})
	// Spawn process 'c1' and 'c2' is actor C on node A for local request from 'a' to 'c1'/'c2'
	nodeA.SpawnRegister("c1", actors.FactoryActorC, gen.ProcessOptions{})
	// nodeA.SpawnRegister("c2", actors.FactoryActorC, gen.ProcessOptions{})

	// Spawm process 'b' is actor B on node B for triggering
	nodeB.SpawnRegister("b", actors.FactoryActorB, gen.ProcessOptions{})
	// Spawn process 'c1' and 'c2' is actor C on node B for local request from 'b' to 'c1'/'c2'
	nodeB.SpawnRegister("c1", actors.FactoryActorC, gen.ProcessOptions{})
	// nodeB.SpawnRegister("c2", actors.FactoryActorC, gen.ProcessOptions{})

	time.Sleep(2 * time.Second)

	// Send trigger message to the process 'a' on node A
	nodeA.Send(gen.Atom("a"), common.DoCallLocal{})
	// Send trigger message to the process 'b' on node B
	// nodeB.Send(gen.Atom("b"), common.DoCallLocal{})

	nodeA.Wait()
}
