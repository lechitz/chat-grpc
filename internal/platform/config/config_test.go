package config

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestConfigValidateSuccess(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			Name:        "chat-grpc",
			Environment: "test",
		},
		Server: ServerConfig{
			Host:           "127.0.0.1",
			Port:           "50051",
			ShutdownGrace:  time.Second,
			MaxRecvMsgSize: 1024,
			MaxSendMsgSize: 1024,
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected validation to succeed, got error: %v", err)
	}
}

func TestConfigValidateFailures(t *testing.T) {
	testCases := []struct {
		name    string
		mutate  func(*Config)
		wantErr error
	}{
		{
			name: "missing app name",
			mutate: func(c *Config) {
				c.App.Name = "   "
			},
			wantErr: ErrAppNameRequired,
		},
		{
			name: "missing environment",
			mutate: func(c *Config) {
				c.App.Environment = ""
			},
			wantErr: ErrEnvironmentRequired,
		},
		{
			name: "missing host",
			mutate: func(c *Config) {
				c.Server.Host = " "
			},
			wantErr: ErrServerHostRequired,
		},
		{
			name: "invalid port text",
			mutate: func(c *Config) {
				c.Server.Port = "invalid"
			},
			wantErr: ErrServerPortInvalid,
		},
		{
			name: "port out of range",
			mutate: func(c *Config) {
				c.Server.Port = "70000"
			},
			wantErr: ErrServerPortInvalid,
		},
		{
			name: "negative shutdown grace",
			mutate: func(c *Config) {
				c.Server.ShutdownGrace = -time.Second
			},
			wantErr: ErrShutdownGraceNegative,
		},
		{
			name: "non positive max recv",
			mutate: func(c *Config) {
				c.Server.MaxRecvMsgSize = 0
			},
			wantErr: ErrMaxRecvSizeInvalid,
		},
		{
			name: "non positive max send",
			mutate: func(c *Config) {
				c.Server.MaxSendMsgSize = -1
			},
			wantErr: ErrMaxSendSizeInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultConfig()
			tc.mutate(cfg)

			err := cfg.Validate()
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestLoadAppliesValidation(t *testing.T) {
	t.Setenv(envPortKey, "70000") // invalid (> 65535)
	resetEnvCache()

	_, err := Load()
	if !errors.Is(err, ErrServerPortInvalid) {
		t.Fatalf("expected ErrServerPortInvalid, got %v", err)
	}
}

func defaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:        "chat-grpc",
			Environment: "dev",
		},
		Server: ServerConfig{
			Host:           "127.0.0.1",
			Port:           "50051",
			ShutdownGrace:  time.Second,
			MaxRecvMsgSize: 1024,
			MaxSendMsgSize: 1024,
		},
	}
}

func resetEnvCache() {
	loadEnvOnce = sync.Once{}
	loadEnvErr = nil
}
