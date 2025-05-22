package types

import (
	"ergo.services/ergo/net/edf"
)

type DoTaskMessage struct {
	TaskId       int64
	TaskProgress int64
	TaskTarget   int64
}

// Register network messages
func init() {
	if err := edf.RegisterTypeOf(DoTaskMessage{}); err != nil {
		panic(err)
	}
}
