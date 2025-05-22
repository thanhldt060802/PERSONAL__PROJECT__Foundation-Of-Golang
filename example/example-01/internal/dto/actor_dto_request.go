package dto

type SimpleRequest struct {
	Body struct {
		Receiver string `json:"receiver" required:"true" doc:"Name of Receiver will send."`
		Message  string `json:"message" required:"true" doc:"Some text for sending to Receiver."`
	}
}

type StartReceiverRequest struct {
	Body struct {
		Receiver  string `json:"receiver" required:"true" doc:"Name of Receiver will start."`
		Task      string `json:"task" required:"true" doc:"Task for sending to Receiver."`
		Situation string `json:"situation" required:"true" doc:"Situation for Receiver to process."`
	}
}
