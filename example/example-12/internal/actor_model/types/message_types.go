package types

type RunTaskMessage struct {
	TaskId int64
}

type GetExistedWorkersMessage struct {
	WorkerNames chan []string
	Running     chan []string
	Available   chan []string
}
