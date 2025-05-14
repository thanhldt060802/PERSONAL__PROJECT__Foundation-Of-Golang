package actors

import (
	"fmt"
	"math/rand"
	"thanhldt060802/common"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ActorC struct {
	act.Actor
}

func FactoryActorC() gen.ProcessBehavior {
	return &ActorC{}
}

func (actorC *ActorC) Init(args ...any) error {
	actorC.Log().Info("started %s process on: %s", actorC.Name(), actorC.Node().Name())
	return nil
}

func (actorC *ActorC) HandleCall(from gen.PID, ref gen.Ref, request any) (any, error) {
	number := rand.Intn(2)

	switch r := request.(type) {
	case common.LocalRequest:
		{
			actorC.Log().Info("received LocalRequest from %s: %#v", from, r)
			// if number == 0 {
			// 	panic("local response error")
			// }
			return fmt.Sprintf("local response for %v", r), nil
		}
	case common.RemoteRequest:
		{
			actorC.Log().Info("received RemoteRequest from %s: %#v", from, r)
			if number == 0 {
				panic("remote response error")
			}
			return fmt.Sprintf("remote response for %v", r), nil
		}
	}

	actorC.Log().Info("received unknown request: %#v", request)
	return nil, nil
}
