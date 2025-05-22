package dto

type SendNewTaskRequest struct {
	Body struct {
		Receiver string `json:"receiver" required:"true" doc:"Name of Receiver will send."`
		TaskId   int64  `json:"task_id" required:"true" doc:"Task id will send to Receiver to process."`
	}
}
