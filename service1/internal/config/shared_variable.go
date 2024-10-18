package config

import(
	"context"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	Ctx             context.Context
	RequestsCounter metric.Int64Counter
	Tracer          trace.Tracer
)