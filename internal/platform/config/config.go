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
