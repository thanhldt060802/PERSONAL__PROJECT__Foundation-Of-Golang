package main

import (
	"fmt"
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
	node1Options.Log.Loggers = append(node1Options.Log.Loggers, gen.Logger{Name: "node_1", Logger: loggerColored})
	node1Options.Network.Cookie = "123"

	node2Options.Log.DefaultLogger.Disable = true
	node2Options.Log.Loggers = append(node2Options.Log.Loggers, gen.Logger{Name: "node_2", Logger: loggerColored})
	node2Options.Network.Cookie = "123"

	node1, err := ergo.StartNode("node1@localhost", node1Options)
	if err != nil {
		panic(err)
	}

	node1.SpawnRegister("sender_1", actors.FactorySenderActor, gen.ProcessOptions{}, &actors.SenderActorParam{ReceiverProcessName: "receiver_1"})
	node1.SpawnRegister("sender_2", actors.FactorySenderActor, gen.ProcessOptions{}, &actors.SenderActorParam{ReceiverProcessName: "receiver_2"})
	node1.SpawnRegister("sender_3", actors.FactorySenderActor, gen.ProcessOptions{}, &actors.SenderActorParam{ReceiverProcessName: "receiver_3"})
	node1.SpawnRegister("sender_4", actors.FactorySenderActor, gen.ProcessOptions{}, &actors.SenderActorParam{ReceiverProcessName: "receiver_4"})
	node1.SpawnRegister("sender_5", actors.FactorySenderActor, gen.ProcessOptions{}, &actors.SenderActorParam{ReceiverProcessName: "receiver_5"})

	node2, err := ergo.StartNode("node2@localhost", node2Options)
	if err != nil {
		panic(err)
	}

	receiverSupervisorPID, _ := node2.Spawn(receiversupervisors.FactoryReceiverSupervisor, gen.ProcessOptions{})
	node2.Log().Info("Supervisor for node2@localhost is started succesfully with PID %s", receiverSupervisorPID)

	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	node1.Send(gen.Atom("sender_1"), "start")
	node1.Send(gen.Atom("sender_2"), "start")
	node1.Send(gen.Atom("sender_3"), "start")
	node1.Send(gen.Atom("sender_4"), "start")
	node1.Send(gen.Atom("sender_5"), "start")

	node1.Wait()
}
