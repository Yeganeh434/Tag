package main

import (
	"context"
	"flag"
	"log"
	"service1/internal/adapters/controllers/http"
	"service1/internal/adapters/databases/mysql"
	"service1/internal/adapters/nats"
	"service1/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/exporters/jaeger"
)

var port = flag.Int("port", 8080, "Port to run the HTTP server")

func main() {
	config.Ctx = context.Background()

	res, err := resource.New(config.Ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("tag_project"),
		),
	)
	if err != nil {
		log.Printf("error creating resource:%v", err)
	}

	metricExporter, err := prometheus.New()
	if err != nil {
		log.Printf("error creating prometheus exporter: %v", err)
	}

	traceExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")),
	)
	
	if err != nil {
		log.Printf("error creating jaeger exporter: %v", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(metricExporter),
		sdkmetric.WithResource(res),
	)
	defer func() {
		if err := meterProvider.Shutdown(config.Ctx); err != nil {
			log.Fatalf("error shutting down meter provider: %v", err)
		}
	}()

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	defer func() {
		if err := tracerProvider.Shutdown(config.Ctx); err != nil {
			log.Fatalf("error shutting down tracer provider: %v", err)
		}
	}()

	otel.SetMeterProvider(meterProvider)
	otel.SetTracerProvider(tracerProvider)

	meter := otel.Meter("gin-metrics")
	config.RequestsCounter, err = meter.Int64Counter(
		"requests_total",
		metric.WithDescription("total number of requests"),
	)
	if err != nil {
		log.Printf("error creating counter: %v", err)
	}

	config.Tracer = otel.Tracer("gin-tracer")

	mysql.InitialDatabase()

	flag.Parse()
	// err := http.RunWebServer()

	go func(){
		nats.Subscribe()
	}()

	err = http.RunWebServer(*port)
	if err != nil {
		log.Printf("could not start server:%v", err)
	}
}
