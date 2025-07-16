package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func StartSpan(ctx context.Context, modulePath string, actionName string) (context.Context, trace.Span) {
	return otel.Tracer(modulePath).Start(ctx, actionName)
}
