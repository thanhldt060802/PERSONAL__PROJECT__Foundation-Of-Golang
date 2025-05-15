package common

import (
	"ergo.services/ergo/net/edf"
)

type LocalRequest struct {
	Message string
}

type RemoteRequest struct {
	Message string
}

// Register network messages
func init() {
	if err := edf.RegisterTypeOf(LocalRequest{}); err != nil {
		panic(err)
	}
	if err := edf.RegisterTypeOf(RemoteRequest{}); err != nil {
		panic(err)
	}
}
