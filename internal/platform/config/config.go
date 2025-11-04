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
	App    AppConfig
	Server ServerConfig
}

// AppConfig holds metadata about the running application.
type AppConfig struct {
	Name        string
	Environment string
}

// ServerConfig hosts gRPC listener configuration.
type ServerConfig struct {
	Host           string
	Port           string
	ShutdownGrace  time.Duration
	MaxRecvMsgSize int
	MaxSendMsgSize int
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
func Load() (*Config, error) {
	if err := loadDotEnv(); err != nil {
		return nil, fmt.Errorf("load env file: %w", err)
	}

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv(envAppNameKey, defaultAppName),
			Environment: getEnv(envEnvironmentKey, defaultEnvironment),
		},
		Server: ServerConfig{
			Host:           getEnv(envHostKey, defaultHost),
			Port:           getEnv(envPortKey, defaultPort),
			ShutdownGrace:  getEnvDuration(envShutdownGraceKey, defaultShutdownGrace),
			MaxRecvMsgSize: getEnvInt(envMaxRecvSizeKey, defaultMaxRecvMsgSize),
			MaxSendMsgSize: getEnvInt(envMaxSendSizeKey, defaultMaxSendMsgSize),
		},
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

// Validation errors returned by Config.Validate.
var (
	ErrAppNameRequired       = errors.New("config: app name is required")
	ErrEnvironmentRequired   = errors.New("config: environment is required")
	ErrServerHostRequired    = errors.New("config: server host is required")
	ErrServerPortInvalid     = errors.New("config: server port must be an integer between 1 and 65535")
	ErrShutdownGraceNegative = errors.New("config: shutdown grace must be zero or positive")
	ErrMaxRecvSizeInvalid    = errors.New("config: max receive message size must be greater than zero")
	ErrMaxSendSizeInvalid    = errors.New("config: max send message size must be greater than zero")
)

// Validate ensures the Config has sane values before it is used by the application.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.App.Name) == "" {
		return ErrAppNameRequired
	}
	if strings.TrimSpace(c.App.Environment) == "" {
		return ErrEnvironmentRequired
	}
	if strings.TrimSpace(c.Server.Host) == "" {
		return ErrServerHostRequired
	}

	port := strings.TrimSpace(c.Server.Port)
	num, err := strconv.Atoi(port)
	if err != nil || num < 1 || num > 65535 {
		return ErrServerPortInvalid
	}

	if c.Server.ShutdownGrace < 0 {
		return ErrShutdownGraceNegative
	}
	if c.Server.MaxRecvMsgSize <= 0 {
		return ErrMaxRecvSizeInvalid
	}
	if c.Server.MaxSendMsgSize <= 0 {
		return ErrMaxSendSizeInvalid
	}

	return nil
}
