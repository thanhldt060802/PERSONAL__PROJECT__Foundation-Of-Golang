package app

import (
	"bytes"
	"encoding/json"
	"net/http"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type MyWebWorker struct {
	act.WebWorker
}

func FactoryMyWebWorker() gen.ProcessBehavior {
	return &MyWebWorker{}
}

// Init invoked on a start this process.
func (myWebWorker *MyWebWorker) Init(args ...any) error {
	myWebWorker.Log().Info("Started WebWorker process with args %v successful", args)
	return nil
}

// Handle GET requests. For the other HTTP methods (POST, PATCH, etc)
// you need to add the accoring callback-method implementation. See act.WebWorkerBehavior.

func (myWebWorker *MyWebWorker) HandleGet(from gen.PID, writer http.ResponseWriter, request *http.Request) error {
	myWebWorker.Log().Info("Got HTTP request %s", request.URL.Path)

	path := request.URL.Path

	data := struct {
		Status  int
		Code    string
		Message string
	}{}

	switch path {
	case "/path1":
		{
			data.Status = http.StatusOK
			data.Code = "OK"
			data.Message = "Hello"
		}
	case "/path2":
		{
			data.Status = http.StatusOK
			data.Code = "OK"
			data.Message = "Goodbye"
		}
	default:
		{
			data.Status = http.StatusNotFound
			data.Code = "ERR_NOT_FOUND"
			data.Message = "Path is not valid"
		}
	}

	var buf bytes.Buffer
	writer.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(data)
	writer.Write(buf.Bytes())

	return nil
}
