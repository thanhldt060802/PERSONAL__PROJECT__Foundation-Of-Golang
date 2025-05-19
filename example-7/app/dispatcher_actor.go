package app

import (
	"math/rand"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type DispatcherActor struct {
	act.Actor
}

func FactoryDispatcherActor() gen.ProcessBehavior {
	return &DispatcherActor{}
}

func (dispatcherActor *DispatcherActor) Init(args ...any) error {
	dispatcherActor.Log().Info("started process %s %s on %s", dispatcherActor.PID(), dispatcherActor.Name(), dispatcherActor.Node().Name())

	dispatcherActor.SendAfter(dispatcherActor.PID(), "start_call", 1*time.Second)

	return nil
}

func (dispatcherActor *DispatcherActor) HandleMessage(from gen.PID, message any) error {
	switch msg := message.(type) {
	case string:
		switch msg {
		case "start_call":
			dispatcherActor.startCallScenario()
		case "next":
			dispatcherActor.runNextStep()
		}
	}
	return nil
}

// Random chọn chiều gọi A -> B hoặc B -> A
func (dispatcherActor *DispatcherActor) startCallScenario() {
	if rand.Intn(2) == 0 {
		dispatcherActor.Log().Info("Simulate: A calls B")
		process := gen.ProcessID{
			Name: gen.Atom("actor1"),
			Node: "node2@localhost",
		}
		dispatcherActor.Send(process, AOffHook)
	} else {
		dispatcherActor.Log().Info("Simulate: B calls A")
		process := gen.ProcessID{
			Name: gen.Atom("actor2"),
			Node: "node2@localhost",
		}
		dispatcherActor.Send(process, CallFromOtherTelephone)
	}

	dispatcherActor.SendAfter(dispatcherActor.PID(), "next", 3*time.Second)
}

// Giả lập tiếp các bước trong cuộc gọi
func (dispatcherActor *DispatcherActor) runNextStep() {
	nextSteps := [][]any{
		{DialledNoBusyOrIncorrect, OwnSideGoesOnHook},
		{BSideAcceptsCall, BSideAnswer, OwnSideGoesOnHook},
		{BSideAcceptsCall, BSideAnswer, OtherSideGoesOnHook, OwnSideGoesOnHook},
	}

	steps := nextSteps[rand.Intn(len(nextSteps))]
	for i, ev := range steps {
		// Gửi xen kẽ vào A hoặc B
		process := gen.ProcessID{
			Name: gen.Atom("actor1"),
			Node: "node2@localhost",
		}
		if i%2 == 1 {
			process = gen.ProcessID{
				Name: gen.Atom("actor2"),
				Node: "node2@localhost",
			}
		}
		dispatcherActor.SendAfter(process, ev, time.Duration(i+1)*2*time.Second)
	}
}
