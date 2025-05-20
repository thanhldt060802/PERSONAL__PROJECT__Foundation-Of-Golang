package types

import (
	"ergo.services/ergo/net/edf"
)

type SimpleMessage struct {
	Message string
}

// Register network messages
func init() {
	if err := edf.RegisterTypeOf(SimpleMessage{}); err != nil {
		panic(err)
	}
}
