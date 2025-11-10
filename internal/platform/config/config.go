// Package config centralizes runtime configuration for the chat-grpc service.
package config

import (
	"bufio"
	"errors"
	"fmt"

	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config is the top-level configuration struct consumed by the application.
type Config struct {
	General       GeneralConfig
	App           AppConfig
	ServerGRPC    ServerConfig
	Observability ObservabilityConfig
}

// AppConfig holds metadata about the running application.
type AppConfig struct {
	Name        string
	Environment string
}

// Addr returns the full host:port pair for binding the gRPC server.
func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

var (
	loadEnvOnce sync.Once
	loadEnvErr  error
)

// Load builds a Config from environment variables, applying sane defaults for local development.
// This is a lightweight wrapper used by tests and callers that don't have a logger available.
func Load() (*Config, error) {
	// reuse existing helpers for reading .env and env getters
	if err := loadDotEnv(); err != nil {
		return nil, fmt.Errorf("load env file: %w", err)
	}

	cfg := &Config{
		General: GeneralConfig{
			Name:    getEnv(envAppNameKey, defaultAppName),
			Env:     getEnv(envEnvironmentKey, defaultEnvironment),
			Version: getEnv(envAppVersionKey, ""),
		},
		App: AppConfig{
			Name:        getEnv(envAppNameKey, defaultAppName),
			Environment: getEnv(envEnvironmentKey, defaultEnvironment),
		},
		ServerGRPC: ServerConfig{
			Host:           getEnv(envHostKey, defaultHost),
			Port:           getEnv(envPortKey, defaultPort),
			ShutdownGrace:  getEnvDuration(envShutdownGraceKey, defaultShutdownGrace),
			MaxRecvMsgSize: getEnvInt(envMaxRecvSizeKey, defaultMaxRecvMsgSize),
			MaxSendMsgSize: getEnvInt(envMaxSendSizeKey, defaultMaxSendMsgSize),
		},
		Observability: ObservabilityConfig{
			Enabled:                  getEnvBool(envOtelEnabledKey, defaultOtelEnabled),
			OtelExporterOTLPEndpoint: getEnv(envOtelEndpointKey, ""),
			OtelExporterInsecure:     getEnvBool(envOtelInsecureKey, defaultOtelInsecure),
			OtelExporterTimeout:      getEnv(envOtelTimeoutKey, defaultOtelTimeout),
			OtelExporterCompression:  getEnv(envOtelCompressionKey, defaultOtelCompression),
			OtelExporterHeaders:      getEnv(envOtelHeadersKey, ""),
			ServiceName:              getEnv(envOtelServiceNameKey, ""),
			ServiceVersion:           getEnv(envOtelServiceVersionKey, defaultOtelServiceVersion),
		},
	}

	if cfg.Observability.ServiceName == "" {
		cfg.Observability.ServiceName = cfg.App.Name
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadDotEnv() error {
	loadEnvOnce.Do(func() {
		loadEnvErr = parseEnvFile(envFileName)
	})
	return loadEnvErr
}

func parseEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		num, err := strconv.Atoi(val)
		if err == nil {
			return num
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		d, err := time.ParseDuration(val)
		if err == nil {
			return d
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		switch strings.ToLower(strings.TrimSpace(val)) {
		case "true", "1", "yes", "y", "on":
			return true
		case "false", "0", "no", "n", "off":
			return false
		}
	}
	return fallback
}

// Validation errors returned by Config.Validate.
var (
	ErrAppNameRequired         = errors.New("config: app name is required")
	ErrEnvironmentRequired     = errors.New("config: environment is required")
	ErrServerHostRequired      = errors.New("config: server host is required")
	ErrServerPortInvalid       = errors.New("config: server port must be an integer between 1 and 65535")
	ErrShutdownGraceNegative   = errors.New("config: shutdown grace must be zero or positive")
	ErrMaxRecvSizeInvalid      = errors.New("config: max receive message size must be greater than zero")
	ErrMaxSendSizeInvalid      = errors.New("config: max send message size must be greater than zero")
	ErrOtelEndpointRequired    = errors.New("config: OTEL exporter endpoint is required when observability is enabled")
	ErrOtelServiceNameRequired = errors.New("config: OTEL service name is required when observability is enabled")
)

// Validate ensures the Config has sane values before it is used by the application.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.App.Name) == "" {
		return ErrAppNameRequired
	}
	if strings.TrimSpace(c.App.Environment) == "" {
		return ErrEnvironmentRequired
	}
	if strings.TrimSpace(c.ServerGRPC.Host) == "" {
		return ErrServerHostRequired
	}

	port := strings.TrimSpace(c.ServerGRPC.Port)
	num, err := strconv.Atoi(port)
	if err != nil || num < 1 || num > 65535 {
		return ErrServerPortInvalid
	}

	if c.ServerGRPC.ShutdownGrace < 0 {
		return ErrShutdownGraceNegative
	}
	if c.ServerGRPC.MaxRecvMsgSize <= 0 {
		return ErrMaxRecvSizeInvalid
	}
	if c.ServerGRPC.MaxSendMsgSize <= 0 {
		return ErrMaxSendSizeInvalid
	}

	if c.Observability.Enabled {
		if strings.TrimSpace(c.Observability.OtelExporterOTLPEndpoint) == "" {
			return ErrOtelEndpointRequired
		}
		if strings.TrimSpace(c.Observability.ServiceName) == "" {
			return ErrOtelServiceNameRequired
		}
	}

	return nil
}
