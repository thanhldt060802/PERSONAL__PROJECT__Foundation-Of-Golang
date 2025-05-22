package dto

type RunTaskRequest struct {
	Body struct {
		WorkerName string `json:"worker_name" required:"true" doc:"Name of Worker will receive task."`
		TaskId     int64  `json:"task_id" required:"true" doc:"Task id will send to Worker to process."`
	}
}

type RunTasksRequest struct {
	Body struct {
		TaskIds []int64 `json:"task_ids" required:"true" doc:"Task id list will distribute to Worker to process."`
	}
}
