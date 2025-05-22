package types

import "ergo.services/ergo/net/edf"

type Task struct {
	TaskId    int64
	Progress  int
	Target    int
	ErrorRate int
}

// Register network messages
func init() {
	if err := edf.RegisterTypeOf(Task{}); err != nil {
		panic(err)
	}
}
