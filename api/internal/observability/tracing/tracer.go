package tracing

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer initializes the OpenTelemetry tracer
// Returns a shutdown function that should be called on application exit
func InitTracer(serviceName, environment string) (func(context.Context) error, error) {
	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
			semconv.DeploymentEnvironment(environment),
		),
	)
	if err != nil {
		return nil, err
	}

	// Determine which exporter to use based on environment
	var exporter sdktrace.SpanExporter

	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint != "" {
		// Use OTLP exporter (for Jaeger, Tempo, etc.)
		log.Printf("Using OTLP exporter with endpoint: %s", otlpEndpoint)
		otlpExporter, err := otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithInsecure(), // Use TLS in production
			otlptracegrpc.WithEndpoint(otlpEndpoint),
		)
		if err != nil {
			return nil, err
		}
		exporter = otlpExporter
	} else {
		// Fallback to stdout exporter (development)
		log.Printf("OTLP endpoint not configured, using stdout exporter")
		stdoutExporter, err := stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			return nil, err
		}
		exporter = stdoutExporter
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Sample all traces (adjust for production)
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator for trace context propagation
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	log.Printf("OpenTelemetry tracer initialized successfully")

	// Return shutdown function
	return tp.Shutdown, nil
}

// Tracer returns the global tracer for the BGC API service
func Tracer() trace.Tracer {
	return otel.Tracer("bgc-api")
}

// StartSpan is a helper function to start a new span
func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return Tracer().Start(ctx, spanName)
}
