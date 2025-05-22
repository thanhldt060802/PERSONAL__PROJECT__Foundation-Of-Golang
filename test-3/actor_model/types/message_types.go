package types

import (
	"ergo.services/ergo/net/edf"
)

type NewTaskMessage struct {
	Receiver string
	TaskId   int64
}

type ExistedReceiverNamesMessage struct {
	ReceiverNames chan []string
}

// Register network messages
func init() {
	if err := edf.RegisterTypeOf(NewTaskMessage{}); err != nil {
		panic(err)
	}
	// if err := edf.RegisterTypeOf(ExistedReceiverNamesMessage{}); err != nil {
	// 	panic(err)
	// }
}
