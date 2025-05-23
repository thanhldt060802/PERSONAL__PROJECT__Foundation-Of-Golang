package app

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/meta"
)

type MyTCP struct {
	act.Actor
}

func FactoryMyTCP() gen.ProcessBehavior {
	return &MyTCP{}
}

func (myTCP *MyTCP) Init(args ...any) error {
	tcpSeverOptions := meta.TCPServerOptions{
		Host: "localhost",
		Port: 12345,
		ProcessPool: []gen.Atom{
			"worker_1",
			"worker_2",
			"worker_3",
		},
	}

	metaTCP, err := meta.CreateTCPServer(tcpSeverOptions)
	if err != nil {
		myTCP.Log().Error("Create TCP server failed: %s", err.Error())
		return err
	}

	id, err := myTCP.SpawnMeta(metaTCP, gen.MetaOptions{})
	if err != nil {
		myTCP.Log().Error("Spawn TCP server meta-process failed: %s", err.Error())
		metaTCP.Terminate(err)
		return err
	}

	myTCP.Log().Info("Started TCP server on %s:%d (meta-process: %s)", tcpSeverOptions.Host, tcpSeverOptions.Port, id)
	myTCP.Log().Info("you may check it with command below:")
	myTCP.Log().Info("   $ ncat %s %d", tcpSeverOptions.Host, tcpSeverOptions.Port)
	return nil
}

// func (myTCP *MyTCP) HandleMessage(from gen.PID, message any) error {
// 	switch receivedMessage := message.(type) {
// 	case meta.MessageTCPConnect:
// 		myTCP.Log().Info("--- New connection with: %s (serving meta-process: %s)", receivedMessage.RemoteAddr, receivedMessage.ID)
// 	case meta.MessageTCPDisconnect:
// 		myTCP.Log().Info("--- Terminated connection (serving meta-process: %s)", receivedMessage.ID)
// 	case meta.MessageTCP:
// 		data := string(receivedMessage.Data)
// 		myTCP.Log().Info("<-- Got TCP packet from %s: %s ", receivedMessage.ID, strings.TrimRight(data, "\r\n"))
// 		receivedMessage.Data = []byte("OK: " + data)
// 		if err := myTCP.SendAlias(receivedMessage.ID, receivedMessage); err != nil {
// 			myTCP.Log().Error("--- Send to %s failed: %s", receivedMessage.ID, err.Error())
// 		}
// 	default:
// 		myTCP.Log().Info("--- Unknown message from %s: %#v", from, message)
// 	}
// 	return nil
// }
