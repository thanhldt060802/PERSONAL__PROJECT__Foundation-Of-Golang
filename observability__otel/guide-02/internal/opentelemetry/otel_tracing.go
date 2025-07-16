package opentelemetry

import (
	"context"
	"fmt"
	"log"
	"thanhldt060802/appconfig"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

type TracerEndPointConfig struct {
	Host string
	Port int
}

func NewTracer(config TracerEndPointConfig) func() {
	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(fmt.Sprintf("%v:%v", config.Host, config.Port)),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err.Error())
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(appconfig.AppConfig.AppName),
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
