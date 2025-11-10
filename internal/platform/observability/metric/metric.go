// Package metric configures the OpenTelemetry meter provider for the chat service.
package metric

import (
	"context"
	"strings"
	"time"

	"github.com/lechitz/chat-grpc/internal/platform/config"
	"github.com/lechitz/chat-grpc/internal/platform/observability"
	"github.com/lechitz/chat-grpc/internal/platform/ports/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

const (
	logWarnInvalidMetricEndpoint = "observability: invalid OTEL exporter endpoint for metrics, using raw value"
	logWarnInvalidMetricTimeout  = "observability: invalid OTEL metric exporter timeout, using default"
	logErrMetricExporterInit     = "observability: failed to initialize OTEL metric exporter"
	logErrMetricShutdown         = "observability: failed to shutdown meter provider"
)

// Init configures the global meter provider when observability is enabled. It returns a cleanup
// function that should be invoked during shutdown to flush metrics.
func Init(cfg *config.Config, log logger.ContextLogger) (func(context.Context), error) {
	if cfg == nil || !cfg.Observability.Enabled {
		return func(context.Context) {}, nil
	}

	endpoint := cfg.Observability.OtelExporterOTLPEndpoint
	exporterEndpoint, err := observability.ExportHostFromEndpoint(endpoint)
	if err != nil {
		log.Warnw(logWarnInvalidMetricEndpoint, "endpoint", endpoint, "error", err)
		exporterEndpoint = endpoint
	}

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(exporterEndpoint),
	}
	if cfg.Observability.OtelExporterInsecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}
	if strings.EqualFold(cfg.Observability.OtelExporterCompression, "gzip") {
		// gRPC exporter currently doesn't support http-style compression option here; ignore or handle separately
	}
	if headers := observability.ParseHeaders(cfg.Observability.OtelExporterHeaders); len(headers) > 0 {
		opts = append(opts, otlpmetricgrpc.WithHeaders(headers))
	}

	// honor timeout by using a context with timeout when creating the exporter
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
			log.Warnw(logWarnInvalidMetricTimeout, "timeout", timeout, "error", err)
		}
	}

	exporter, err := otlpmetricgrpc.New(expCtx, opts...)
	if err != nil {
		log.Errorw(logErrMetricExporterInit, "error", err)
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.Observability.ServiceName),
		semconv.ServiceVersionKey.String(cfg.Observability.ServiceVersion),
	)

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	return func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := provider.Shutdown(ctx); err != nil {
			log.Errorw(logErrMetricShutdown, "error", err)
		}
	}, nil
}
