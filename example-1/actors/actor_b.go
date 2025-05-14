package actors

import (
	"fmt"
	"math/rand"
	"thanhldt060802/common"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ActorB struct {
	act.Actor
	count int
}

func FactoryActorB() gen.ProcessBehavior {
	return &ActorB{}
}

func (actorB *ActorB) Init(args ...any) error {
	actorB.Log().Info("started %s process on: %s", actorB.Name(), actorB.Node().Name())
	actorB.count = 1
	return nil
}

func (actorB *ActorB) HandleMessage(from gen.PID, message any) error {
	actorName := "c1"
	delayTime := time.Duration(rand.Intn(int(200*time.Millisecond-100*time.Millisecond))) + 100*time.Millisecond

	switch message.(type) {
	case common.DoCallLocal:
		{
			process := gen.Atom(actorName)
			actorB.Log().Info("making request to local process %s", process)
			if result, err := actorB.Call(process, common.LocalRequest{Message: fmt.Sprintf("Message %d of process 'b'", actorB.count)}); err == nil {
				actorB.Log().Info("received result from local process %s: %#v", process, result)
			} else {
				actorB.Log().Error("call local process %s failed: %s", process, err.Error())
			}
			actorB.count++
			actorB.SendAfter(gen.Atom("b"), common.DoCallRemote{}, delayTime)
			return nil
		}
	case common.DoCallRemote:
		{
			process := gen.ProcessID{Name: gen.Atom(actorName), Node: "node2@localhost"}
			actorB.Log().Info("making request to remote process %s", process.Name)
			if result, err := actorB.Call(process, common.RemoteRequest{Message: fmt.Sprintf("Message %d of process 'b'", actorB.count)}); err == nil {
				actorB.Log().Info("received result from remote process %s: %#v", process.Name, result)
			} else {
				actorB.Log().Error("call remote process %s failed: %s", process.Name, err.Error())
			}
			actorB.count++
			actorB.SendAfter(gen.Atom("b"), common.DoCallLocal{}, delayTime)
			return nil
		}
	}

	actorB.Log().Error("unknown message %#v", message)
	return nil
}
