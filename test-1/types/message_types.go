package types

import (
	"ergo.services/ergo/net/edf"
)

type SimpleMessage struct {
	Message string
}

type TaskMessage struct {
	Receiver  string
	Task      string
	Situation string
}

type Task struct {
	Task      string
	Situation string
}

// Register network messages
func init() {
	if err := edf.RegisterTypeOf(SimpleMessage{}); err != nil {
		panic(err)
	}
}
