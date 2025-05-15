package main

import (
	"math/rand"
	"thanhldt060802/actors"
	"thanhldt060802/infrastructure"
	receiversupervisors "thanhldt060802/receiver_supervisors"
	"thanhldt060802/repository"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	infrastructure.InitPostgesDB()
	defer infrastructure.PostgresDB.Close()
	repository.InitTaskRepository()

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

	node1.SpawnRegister("sender-1", actors.FactorySenderActor, gen.ProcessOptions{}, "receiver-1")
	node1.SpawnRegister("sender-2", actors.FactorySenderActor, gen.ProcessOptions{}, "receiver-2")
	node1.SpawnRegister("sender-3", actors.FactorySenderActor, gen.ProcessOptions{}, "receiver-3")
	node1.SpawnRegister("sender-4", actors.FactorySenderActor, gen.ProcessOptions{}, "receiver-4")
	node1.SpawnRegister("sender-5", actors.FactorySenderActor, gen.ProcessOptions{}, "receiver-5")

	node2, err := ergo.StartNode("node2@localhost", node2Options)
	if err != nil {
		panic(err)
	}

	receiverSupervisorPID, _ := node2.Spawn(receiversupervisors.FactoryReceiverSupervisor, gen.ProcessOptions{})
	node2.Log().Info("Supervisor for node 2 is started succesfully with PID %s", receiverSupervisorPID)

	time.Sleep(2 * time.Second)

	node1.Send(gen.Atom("sender-1"), "start")
	node1.Send(gen.Atom("sender-2"), "start")
	node1.Send(gen.Atom("sender-3"), "start")
	node1.Send(gen.Atom("sender-4"), "start")
	node1.Send(gen.Atom("sender-5"), "start")

	node1.Wait()
}
