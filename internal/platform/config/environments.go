package config

import "time"

// GeneralConfig holds a general application configuration.
type GeneralConfig struct {
	Name    string `envconfig:"APP_NAME"`
	Env     string `envconfig:"APP_ENV"`
	Version string `envconfig:"APP_VERSION"`
}

// ObservabilityConfig holds all observability-related configuration.
type ObservabilityConfig struct {
	Enabled                  bool   `envconfig:"OTEL_ENABLED"                default:"false"`
	OtelExporterOTLPEndpoint string `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT" default:"otel-collector:4318"`
	OtelServiceName          string `envconfig:"OTEL_SERVICE_NAME"           default:"AionApi"`
	OtelServiceVersion       string `envconfig:"OTEL_SERVICE_VERSION"        default:"0.1.0"`
	OtelExporterHeaders      string `envconfig:"OTEL_EXPORTER_HEADERS"       default:""`
	OtelExporterTimeout      string `envconfig:"OTEL_EXPORTER_TIMEOUT"       default:"5s"`
	OtelExporterCompression  string `envconfig:"OTEL_EXPORTER_COMPRESSION"   default:"none"`
	OtelExporterInsecure     bool   `envconfig:"OTEL_EXPORTER_INSECURE"      default:"true"`
	ServiceName              string `envconfig:"SERVICE_NAME"`
	ServiceVersion           string `envconfig:"SERVICE_VERSION"`
}

// ServerConfig hosts gRPC listener configuration.
type ServerConfig struct {
	Host           string
	Port           string
	ShutdownGrace  time.Duration
	MaxRecvMsgSize int
	MaxSendMsgSize int
}
