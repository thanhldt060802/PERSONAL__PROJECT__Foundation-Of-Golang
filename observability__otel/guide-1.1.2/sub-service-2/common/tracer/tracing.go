package tracer

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// INTERNAL USAGE

func StartSpanInternal(ctx context.Context, modulePath string, actionName string) (context.Context, trace.Span) {
	return otel.Tracer(modulePath).Start(ctx, actionName)
}

// CROSS SERVICE USAGE

func StartSpanCrossService(ctx context.Context, modulePath string, actionName string, method string, url string) (context.Context, trace.Span, *http.Request, error) {
	ctx, span := otel.Tracer(modulePath).Start(ctx, actionName)

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, nil, nil, err
	}

	return ctx, span, req, nil
}

// PUB/SUB USAGE

type MessageTracing struct {
	TraceContext propagation.MapCarrier
	Payload      any
}

func (msgTrace *MessageTracing) ExtractTraceContext() context.Context {
	return otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(msgTrace.TraceContext))
}
