// Package tracer configures the OpenTelemetry tracer provider for the chat service.
package tracer

import (
	"context"
	"strings"
	"time"

	"github.com/lechitz/chat-grpc/internal/platform/config"
	"github.com/lechitz/chat-grpc/internal/platform/observability"
	"github.com/lechitz/chat-grpc/internal/platform/ports/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

const (
	logWarnInvalidEndpoint = "observability: invalid OTEL exporter endpoint, using raw value"
	logWarnInvalidTimeout  = "observability: invalid OTEL exporter timeout, using default"
	logErrExporterInit     = "observability: failed to initialize OTEL trace exporter"
	logErrShutdownTracer   = "observability: failed to shutdown tracer provider"
)

// Init configures the global tracer provider when observability is enabled. It returns a cleanup
// function that should be invoked during shutdown to flush spans.
func Init(cfg *config.Config, log logger.ContextLogger) (func(context.Context), error) {
	if cfg == nil || !cfg.Observability.Enabled {
		return func(context.Context) {}, nil
	}

	endpoint := cfg.Observability.OtelExporterOTLPEndpoint
	exporterEndpoint, err := observability.ExportHostFromEndpoint(endpoint)
	if err != nil {
		log.Warnw(logWarnInvalidEndpoint, "endpoint", endpoint, "error", err)
		exporterEndpoint = endpoint
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(exporterEndpoint),
	}
	if cfg.Observability.OtelExporterInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	if headers := observability.ParseHeaders(cfg.Observability.OtelExporterHeaders); len(headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(headers))
	}

	expCtx := context.Background()
	var cancel func()
	if timeout := strings.TrimSpace(cfg.Observability.OtelExporterTimeout); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			expCtx, cancel = context.WithTimeout(context.Background(), d)
			defer func() {
				if cancel != nil {
					cancel()
				}
			}()
		} else {
			log.Warnw(logWarnInvalidTimeout, "timeout", timeout, "error", err)
		}
	}

	exporter, err := otlptracegrpc.New(expCtx, opts...)
	if err != nil {
		log.Errorw(logErrExporterInit, "error", err)
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.Observability.ServiceName),
		semconv.ServiceVersionKey.String(cfg.Observability.ServiceVersion),
	)

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(provider)

	return func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := provider.Shutdown(ctx); err != nil {
			log.Errorw(logErrShutdownTracer, "error", err)
		}
	}, nil
}
