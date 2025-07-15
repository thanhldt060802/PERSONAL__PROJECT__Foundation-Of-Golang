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

type apiUser struct {
	tracer      trace.Tracer
	userService service.IUserService
}

func RegisterAPIExample(api hureg.APIGen, userService service.IUserService) {
	handler := &apiUser{
		tracer:      otel.Tracer("api/v1/user.go"),
		userService: userService,
	}

	apiGroup := api.AddBasePath("/user")

	hureg.Register(
		apiGroup,
		huma.Operation{
			OperationID: "user-get-list",
			Method:      http.MethodGet,
			Path:        "",
			Security:    authMdw.DefaultAuthSecurity,
			Description: "Get list users.",
			Middlewares: huma.Middlewares{authMdw.NewAuthMiddleware(api)},
		},
		handler.Gets,
	)
}

type GetsUserResponse struct {
	Body struct {
		Data []*model.User `json:"data" doc:"List user data"`
	}
}

func (handler *apiUser) Gets(ctx context.Context, req *struct{}) (res *GetsUserResponse, err error) {
	ctx, span1 := handler.tracer.Start(ctx, "Handler Gets()")
	defer span1.End()

	data, err := handler.userService.Gets(ctx)
	if err != nil {
		span1.SetStatus(codes.Error, err.Error())
		return
	}

	span1.SetAttributes(attribute.Int("total", len(data)))
	span1.SetStatus(codes.Ok, "success")
	res = &GetsUserResponse{}
	res.Body.Data = data
	return
}
