package otelclient

import (
	"context"
	"log"
	server "thanhldt060802/server/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func InitTracer() func() {
	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("localhost:4318"), // Jaeger HTTP OTLP
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err.Error())
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(server.SERVICE_NAME),
	)

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)

	return func() {
		_ = tracerProvider.Shutdown(ctx)
	}
}
