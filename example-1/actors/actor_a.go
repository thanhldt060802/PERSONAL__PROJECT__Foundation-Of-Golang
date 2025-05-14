package actors

import (
	"fmt"
	"math/rand"
	"thanhldt060802/common"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ActorA struct {
	act.Actor
	count int
}

func FactoryActorA() gen.ProcessBehavior {
	return &ActorA{}
}

func (actorA *ActorA) Init(args ...any) error {
	actorA.Log().Info("started %s process on: %s", actorA.Name(), actorA.Node().Name())
	actorA.count = 1
	return nil
}

func (actorA *ActorA) HandleMessage(from gen.PID, message any) error {
	actorName := "c1"
	delayTime := time.Duration(rand.Intn(int(200*time.Millisecond-100*time.Millisecond))) + 100*time.Millisecond

	switch message.(type) {
	case common.DoCallLocal:
		{
			process := gen.Atom(actorName)
			actorA.Log().Info("making request to local process %s", process)
			if result, err := actorA.Call(process, common.LocalRequest{Message: fmt.Sprintf("Message %d of processs 'a'", actorA.count)}); err == nil {
				actorA.Log().Info("received result from local process %s: %#v", process, result)
			} else {
				actorA.Log().Error("call local process %s failed: %s", process, err.Error())
			}
			actorA.count++
			actorA.SendAfter(gen.Atom("a"), common.DoCallRemote{}, delayTime)
			return nil
		}
	case common.DoCallRemote:
		{
			process := gen.ProcessID{Name: gen.Atom(actorName), Node: "node2@localhost"}
			actorA.Log().Info("making request to remote process %s", process.Name)
			if result, err := actorA.Call(process, common.RemoteRequest{Message: fmt.Sprintf("Message %d of process 'a'", actorA.count)}); err == nil {
				actorA.Log().Info("received result from remote process %s: %#v", process.Name, result)
			} else {
				actorA.Log().Error("call remote process %s failed: %s", process.Name, err.Error())
			}
			actorA.count++
			actorA.SendAfter(gen.Atom("a"), common.DoCallLocal{}, delayTime)
			return nil
		}
	}

	actorA.Log().Error("unknown message %#v", message)
	return nil
}
