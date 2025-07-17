package v1

import (
	"context"
	"net/http"
	"thanhldt060802/common/tracer"
	"thanhldt060802/service"

	authMdw "thanhldt060802/middleware/auth"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
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
			OperationID: "player-get-by-id",
			Method:      http.MethodGet,
			Path:        "/{player_uuid}",
			Security:    authMdw.DefaultAuthSecurity,
			Description: "Get player by id.",
			Middlewares: huma.Middlewares{authMdw.NewAuthMiddleware(api)},
		},
		handler.GetById,
	)
}

type GetPlayerByIdResponse struct {
	Body struct {
		Result string `json:"result" doc:"Result"`
	}
}

func (handler *apiPlayer) GetById(ctx context.Context, req *struct {
	PlayerUuid string `path:"player_uuid" format:"uuid" doc:"Player uuid"`
}) (res *GetPlayerByIdResponse, err error) {
	ctx, span := tracer.StartSpanInternal(ctx, "api/v1/player.go", "Handler.GetById")
	defer span.End()

	result, err := handler.playerService.GetById(ctx, req.PlayerUuid)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	res = &GetPlayerByIdResponse{}
	res.Body.Result = result
	span.SetStatus(codes.Ok, "success")
	return
}
