package dto

import (
	"ergo.services/ergo/net/edf"
)

type TaskRequest struct {
	Id int64
}

// Register network messages
func init() {
	if err := edf.RegisterTypeOf(TaskRequest{}); err != nil {
		panic(err)
	}
}
