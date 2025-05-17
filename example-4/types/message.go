package types

import (
	"ergo.services/ergo/net/edf"
)

type DoTaskMessage struct {
	Id       int64
	Progress int64
	Target   int64
}

type ReturnTaskMessage struct {
	Id         int64
	Progress   int64
	Target     int64
	SelfStatus string
}

// Register network messages
func init() {
	if err := edf.RegisterTypeOf(DoTaskMessage{}); err != nil {
		panic(err)
	}
	if err := edf.RegisterTypeOf(ReturnTaskMessage{}); err != nil {
		panic(err)
	}
}
