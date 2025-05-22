package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	actormodelapp "thanhldt060802/actor_model_app"
	"thanhldt060802/internal/dto"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

var receiverNode gen.Node

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Set options for node1 and node2

	// var node1Options, node2Options gen.NodeOptions
	var node2Options gen.NodeOptions

	optionColored := colored.Options{TimeFormat: time.DateTime, IncludeName: true}
	loggerColored, err := colored.CreateLogger(optionColored)
	if err != nil {
		panic(err)
	}

	// node1Options.Log.DefaultLogger.Disable = true
	// node1Options.Log.Loggers = append(node1Options.Log.Loggers, gen.Logger{Name: "node1", Logger: loggerColored})
	// node1Options.Network.Cookie = "123"

	node2Options.Log.DefaultLogger.Disable = true
	node2Options.Log.Loggers = append(node2Options.Log.Loggers, gen.Logger{Name: "node2", Logger: loggerColored})
	node2Options.Network.Cookie = "123"

	// Init node1 which owns sender_1, sender_2 and local receiver_1

	// node1, err := ergo.StartNode(gen.Atom("node1@localhost"), node1Options)
	// if err != nil {
	// 	panic(err)
	// }

	// node1.SpawnRegister(gen.Atom("sender"), actormodelapp.FactorySenderActor, gen.ProcessOptions{})

	// senderNode = node1

	// Init node2 which owns receiver_1

	node2, err := ergo.StartNode(gen.Atom("node2@localhost"), node2Options)
	if err != nil {
		panic(err)
	}

	node2.SpawnRegister(gen.Atom("receiver"), actormodelapp.FactoryReceiverFSMActor, gen.ProcessOptions{})

	receiverNode = node2

	fmt.Println()
	fmt.Println()

	// Huma Docs UI template by Scalar
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
		Path:        "/fsm",
		Summary:     "/fsm",
		Description: "Get next state from event.",
		Tags:        []string{"User"},
	}, GetNextState)

	r.Run(":8080")

}

func GetNextState(ctx context.Context, reqDTO *dto.SimpleRequest) (*dto.BodyResponse[string], error) {
	nextState := make(chan string)
	message := dto.MyRequest{
		Event:     reqDTO.Body.Event,
		NextState: nextState,
	}
	receiverNode.Send(gen.Atom("receiver"), message)

	select {
	case result := <-nextState:
		{
			res := &dto.BodyResponse[string]{}
			res.Body.Code = "OK"
			res.Body.Message = "Get next state from event successful"
			res.Body.Data = result
			return res, nil
		}
	case <-time.After(3 * time.Second):
		{
			res := &dto.ErrorResponse{}
			res.Status = http.StatusRequestTimeout
			res.Code = "ERR_REQUEST_TIMEOUT"
			res.Message = "Get next state from event failed"
			res.Details = []string{"request time out"}
			return nil, res
		}
	}
}
