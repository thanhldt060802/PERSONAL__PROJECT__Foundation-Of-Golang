package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var DefaultAuthSecurity = []map[string][]string{
	{"standard-auth": {""}},
}

type IAuthMiddleware interface {
	AuthMiddleware(ctx context.Context) (string, error)
}

var AuthMdw IAuthMiddleware
var tracer trace.Tracer = otel.Tracer("middleware/auth.go")

func NewAuthMiddleware(api hureg.APIGen) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		log.Info("========> standard-auth middelware request")
		isAuthorizationRequired := false
		for _, opScheme := range ctx.Operation().Security {
			var ok bool
			if _, ok = opScheme["standard-auth"]; ok {
				log.Info("========> standard-auth middelware validate")
				isAuthorizationRequired = true
				break
			}
		}
		log.Infof("========> require authorization: %v", isAuthorizationRequired)
		if isAuthorizationRequired {
			HumaAuthMiddleware(api, ctx, next)
		} else {
			next(ctx)
		}
	}
}

func HumaAuthMiddleware(api hureg.APIGen, ctx huma.Context, next func(huma.Context)) {
	standardCtx, span1 := tracer.Start(ctx.Context(), "Middleware HumaAuthMiddleware()")
	defer span1.End()

	authHeaderValue := ctx.Header("Authorization")
	if len(authHeaderValue) < 1 {
		log.Error("========> invalid credentials")
		huma.WriteErr(api.GetHumaAPI(), ctx, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), errors.New("missing token"))
		span1.SetStatus(codes.Error, "missing token")
		return
	}

	ctx = huma.WithContext(ctx, standardCtx)
	ctx = huma.WithValue(ctx, "auth-header", authHeaderValue)

	authResult, err := AuthMdw.AuthMiddleware(ctx.Context())
	if err != nil {
		huma.WriteErr(api.GetHumaAPI(), ctx, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
		span1.SetStatus(codes.Error, err.Error())
		return
	}
	log.Infof("========> auth result: %v", authResult)

	ctx = huma.WithContext(ctx, ctx.Context())

	span1.SetStatus(codes.Ok, "success")
	next(ctx)
}
