package types

type RunTaskMessage struct {
	WorkerName string
	TaskId     int64
}

type RunTasksMessage struct {
	TaskIds []int64
}

type GetExistedWorkersMessage struct {
	WorkerNames chan []string
}

// Register network messages
func init() {
	// if err := edf.RegisterTypeOf(RunTaskMessage{}); err != nil {
	// 	panic(err)
	// }
}
