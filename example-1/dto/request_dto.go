package dto

import (
	"ergo.services/ergo/net/edf"
)

type SimpleRequest struct {
	Message string
}

// Register network messages
func init() {
	if err := edf.RegisterTypeOf(SimpleRequest{}); err != nil {
		panic(err)
	}
}
