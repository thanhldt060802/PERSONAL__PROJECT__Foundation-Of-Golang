package v1

import (
	"context"
	"net/http"
	"thanhldt060802/model"
	"thanhldt060802/service"

	authMdw "thanhldt060802/middleware/auth"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type apiPlayer struct {
	tracer        trace.Tracer
	playerService service.IPlayerService
}

func RegisterAPIExample(api hureg.APIGen, playerService service.IPlayerService) {
	handler := &apiPlayer{
		tracer:        otel.Tracer("api/v1/player.go"),
		playerService: playerService,
	}

	apiGroup := api.AddBasePath("/player")

	hureg.Register(
		apiGroup,
		huma.Operation{
			OperationID: "player-get-list",
			Method:      http.MethodGet,
			Path:        "",
			Security:    authMdw.DefaultAuthSecurity,
			Description: "Get list players.",
			Middlewares: huma.Middlewares{authMdw.NewAuthMiddleware(api)},
		},
		handler.Get,
	)
}

type GetsPlayerResponse struct {
	Body struct {
		Data []*model.Player `json:"data" doc:"List player data"`
	}
}

func (handler *apiPlayer) Get(ctx context.Context, req *struct{}) (res *GetsPlayerResponse, err error) {
	ctx, span1 := handler.tracer.Start(ctx, "Handler Get()")
	defer span1.End()

	data, err := handler.playerService.Get(ctx)
	if err != nil {
		span1.SetStatus(codes.Error, err.Error())
		return
	}

	span1.SetAttributes(attribute.Int("total", len(data)))
	span1.SetStatus(codes.Ok, "success")
	res = &GetsPlayerResponse{}
	res.Body.Data = data
	return
}
