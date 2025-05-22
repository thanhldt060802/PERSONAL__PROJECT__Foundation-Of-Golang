package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"thanhtldt060802/app"
	"thanhtldt060802/internal/dto"
	"thanhtldt060802/types"
	"time"

	"ergo.services/application/observer"
	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

var receiverNode gen.Node
var supervisorPID gen.PID

func main() {

	rand.New(rand.NewSource(time.Now().UnixNano()))

	var nodeOptions gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	nodeOptions.Applications = []gen.ApplicationBehavior{
		observer.CreateApp(observer.Options{}),
	}
	nodeOptions.Log.DefaultLogger.Disable = true
	nodeOptions.Log.Loggers = append(nodeOptions.Log.Loggers, gen.Logger{Name: "my-node", Logger: loggerColored})
	nodeOptions.Network.Cookie = "123"

	myNode, err := ergo.StartNode(gen.Atom("mynode@localhost"), nodeOptions)
	if err != nil {
		panic(err)
	}

	receiverNode = myNode

	supervisorPID, _ = myNode.Spawn(app.FactoryReceiverSupervisor, gen.ProcessOptions{})

	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	myNode.Send(supervisorPID, nil)

	// select {}

	var humaDocsEmbedded = `<!doctype html>
	<html>
	  <head>
	    <title>FashionECom APIs</title>
	    <meta charset="utf-8" />
	    <meta name="viewport" content="width=device-width, initial-scale=1" />
	  </head>
	  <body>
	    <script
	      id="api-reference"
	      data-url="/openapi.json"></script>
	    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
	  </body>
	</html>`

	humaCfg := huma.DefaultConfig("Test FSM Actor", "v1.0.0")
	humaCfg.DocsPath = ""
	humaCfg.JSONSchemaDialect = ""
	humaCfg.CreateHooks = nil

	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		details := make([]string, len(errs))
		for i, err := range errs {
			details[i] = err.Error()
		}
		res := &dto.ErrorResponse{}
		res.Status = status
		res.Message = msg
		res.Details = details
		return res
	}

	r := gin.Default()
	r.GET("/docs", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html", []byte(humaDocsEmbedded))
	})

	api := humagin.New(r, humaCfg)

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/send-to-receiver",
		Summary:     "/send-to-receiver",
		Description: "Send simple request to receiver.",
		Tags:        []string{"Test"},
	}, SendSimpleRequest)

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/start-receiver",
		Summary:     "/start-receiver",
		Description: "Send start receiver request to supervisor.",
		Tags:        []string{"Test"},
	}, StartReceiverRequest)

	r.Run(":8080")

}

func SendSimpleRequest(ctx context.Context, reqDTO *dto.SimpleRequest) (*dto.SuccessResponse, error) {
	if err := receiverNode.Send(gen.Atom(reqDTO.Body.Receiver), reqDTO.Body.Message); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Send simple request to receiver failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Send simple request to receiver successful"
	return res, nil
}

func StartReceiverRequest(ctx context.Context, reqDTO *dto.StartReceiverRequest) (*dto.SuccessResponse, error) {
	if err := receiverNode.Send(supervisorPID, types.TaskMessage{Receiver: reqDTO.Body.Receiver, Task: reqDTO.Body.Task, Situation: reqDTO.Body.Situation}); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Send start receiver request to supervisor failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Send start receiver request to supervisor successful"
	return res, nil
}
