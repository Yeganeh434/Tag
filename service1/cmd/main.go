package main

import (
	"context"
	"flag"
	"log"
	"tag_project/internal/adapters/controllers/http"
	"tag_project/internal/adapters/databases/mysql"
	"tag_project/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	// jaeger "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	// "go.opentelemetry.io/otel/trace"

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
	// traceExporter,err:=jaeger.New(context.Background(),
	// 	jaeger.WithEndpoint("http://localhost:14268/api/traces"),
	// )
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

	err = http.RunWebServer(*port)
	if err != nil {
		log.Printf("could not start server:%v", err)
	}
}

// package main

// import (
// 	"context"
// 	"flag"
// 	"log"
// 	"tag_project/internal/adapters/controllers/http"
// 	"tag_project/internal/adapters/databases/mysql"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/exporters/prometheus"
// 	sdkmetric "go.opentelemetry.io/otel/sdk/metric" // تغییر نام این پکیج به sdkmetric
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// 	jaeger "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
// 	"go.opentelemetry.io/otel/sdk/resource"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
// 	"go.opentelemetry.io/otel/metric"
// 	"go.opentelemetry.io/otel/trace"
// )

// var (
// 	port            = flag.Int("port", 8080, "Port to run the HTTP server")
// 	Ctx             context.Context
// 	RequestsCounter metric.Int64Counter
// 	Tracer          trace.Tracer
// )

// func main() {
// 	Ctx = context.Background()

// 	res, err := resource.New(Ctx,
// 		resource.WithAttributes(
// 			semconv.ServiceNameKey.String("tag_project"),
// 		),
// 	)
// 	if err != nil {
// 		log.Printf("error creating resource: %v", err)
// 	}

// 	metricExporter, err := prometheus.New()
// 	if err != nil {
// 		log.Printf("error creating prometheus exporter: %v", err)
// 	}

// 	traceExporter, err := jaeger.New(context.Background(),
// 		jaeger.WithEndpoint("http://localhost:14268/api/traces"),
// 		jaeger.WithInsecure(),
// 	)
// 	if err != nil {
// 		log.Printf("error creating jaeger exporter: %v", err)
// 	}

// 	// استفاده از sdkmetric به جای metric برای ایجاد MeterProvider
// 	meterProvider := sdkmetric.NewMeterProvider(
// 		sdkmetric.WithReader(metricExporter),
// 		sdkmetric.WithResource(res),
// 	)
// 	defer func() {
// 		if err := meterProvider.Shutdown(Ctx); err != nil {
// 			log.Fatalf("error shutting down meter provider: %v", err)
// 		}
// 	}()

// 	tracerProvider := sdktrace.NewTracerProvider(
// 		sdktrace.WithBatcher(traceExporter),
// 		sdktrace.WithResource(res),
// 	)
// 	defer func() {
// 		if err := tracerProvider.Shutdown(Ctx); err != nil {
// 			log.Fatalf("error shutting down tracer provider: %v", err)
// 		}
// 	}()

// 	otel.SetMeterProvider(meterProvider)
// 	otel.SetTracerProvider(tracerProvider)

// 	// استفاده از API جدید برای ایجاد Counter
// 	meter := otel.Meter("gin-metrics")
// 	RequestsCounter, err = meter.Int64Counter(
// 		"requests_total",
// 		metric.WithDescription("Total number of requests"),
// 	)
// 	if err != nil {
// 		log.Printf("error creating counter: %v", err)
// 	}

// 	Tracer = otel.Tracer("gin-tracer")

// 	mysql.InitialDatabase()

// 	flag.Parse()
// 	err = http.RunWebServer(*port)
// 	if err != nil {
// 		log.Printf("could not start server: %v", err)
// 	}
// }
