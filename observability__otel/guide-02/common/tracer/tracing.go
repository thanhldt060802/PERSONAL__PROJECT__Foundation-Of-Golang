package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type TracingManager struct {
	tracingMap map[string]trace.Tracer
}

var TracingMg *TracingManager

func NewTracingManager() *TracingManager {
	return &TracingManager{
		tracingMap: map[string]trace.Tracer{},
	}
}

func (tracingManager *TracingManager) NewSpan(ctx context.Context, location string, spanName string) (context.Context, trace.Span) {
	if _, ok := tracingManager.tracingMap[location]; !ok {
		tracingManager.tracingMap[location] = otel.Tracer(location)
	}
	return tracingManager.tracingMap[location].Start(ctx, spanName)
}
