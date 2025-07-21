package main

import (
	"context"
	"errors"
	"math/rand"
	"thanhldt060802/common/tracer"
	"thanhldt060802/internal/opentelemetry"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var EXAMPLE_NUM int = 1
var EXAMPLES map[int]func()

func init() {
	EXAMPLES = map[int]func(){
		1: Example1,
		2: Example2,
		3: Example3,
	}
}

func main() {

	EXAMPLES[EXAMPLE_NUM]()

}

// Example for internal tracing.
func Example1() {
	opentelemetry.ShutdownTracer = opentelemetry.NewTracer(opentelemetry.TracerEndPointConfig{
		ServiceName: "my-service",
		Host:        "localhost",
		Port:        4318,
	})
	defer opentelemetry.ShutdownTracer()

	function3 := func(ctx context.Context) {
		_, span := tracer.StartSpanInternal(ctx)
		defer span.End()

		span.AddEvent("Fetch data", trace.WithAttributes(
			attribute.String("data", "Some data"),
		))
		time.Sleep(1 * time.Second)

		span.AddEvent("Fetch extensive data")
		time.Sleep(1 * time.Second)
	}
	function2 := func(ctx context.Context) {
		ctx, span := tracer.StartSpanInternal(ctx)
		defer span.End()

		span.AddEvent("Call to function3", trace.WithAttributes(
			attribute.String("param.a", "Some data"),
			attribute.String("param.b", "Some data"),
		))
		function3(ctx)
		time.Sleep(1 * time.Second)

		if rand.Intn(2) == 0 {
			span.Err = errors.New("some thing wrong")
		}
	}
	function1 := func(ctx context.Context) {
		ctx, span := tracer.StartSpanInternal(ctx)
		defer span.End()

		span.AddEvent("Call to function2")
		function2(ctx)
		time.Sleep(1 * time.Second)

		span.SetAttributes(
			attribute.String("result", "Some data"),
			attribute.String("message", "Some data"),
		)
	}

	function1(context.Background())

	select {}
}

// Example for cross-service tracing.
func Example2() {

}

// Example for pub/sub system tracing.
func Example3() {

}
