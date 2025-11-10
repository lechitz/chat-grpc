package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/chat-grpc/internal/platform/ports/logger"
	"github.com/lechitz/chat-grpc/internal/shared/commonkeys"
)

// Loader is responsible for reading environment configuration.
type Loader struct {
	cfg Config
}

// New returns a new instance of Loader.
func New() *Loader {
	return &Loader{}
}

// Load reads environment configuration and returns a Config struct.
// Returns an error if environment parsing or key generation fails.
func (l *Loader) Load(logger logger.ContextLogger) (*Config, error) {
	if err := envconfig.Process(commonkeys.Setting, &l.cfg); err != nil {
		logger.Errorw(ErrFailedToProcessEnvVars, commonkeys.Error, err)
		return nil, err
	}

	if l.cfg.General.Name == "" {
		l.cfg.General.Name = getEnv(envAppNameKey, defaultAppName)
	}
	if l.cfg.General.Env == "" {
		l.cfg.General.Env = getEnv(envEnvironmentKey, defaultEnvironment)
	}

	if l.cfg.General.Version == "" {
		l.cfg.General.Version = getEnv("CHAT_GRPC_APP_VERSION", "")
	}

	if l.cfg.App.Name == "" {
		l.cfg.App.Name = defaultAppName
	}
	if l.cfg.App.Environment == "" {
		l.cfg.App.Environment = defaultEnvironment
	}
	if l.cfg.ServerGRPC.Host == "" {
		l.cfg.ServerGRPC.Host = getEnv(envHostKey, defaultHost)
	}
	if l.cfg.ServerGRPC.Port == "" {
		l.cfg.ServerGRPC.Port = getEnv(envPortKey, defaultPort)
	}
	if l.cfg.ServerGRPC.ShutdownGrace == 0 {
		l.cfg.ServerGRPC.ShutdownGrace = defaultShutdownGrace
	}
	if l.cfg.ServerGRPC.MaxRecvMsgSize == 0 {
		l.cfg.ServerGRPC.MaxRecvMsgSize = defaultMaxRecvMsgSize
	}
	if l.cfg.ServerGRPC.MaxSendMsgSize == 0 {
		l.cfg.ServerGRPC.MaxSendMsgSize = defaultMaxSendMsgSize
	}

	if l.cfg.Observability.ServiceName == "" {
		l.cfg.Observability.ServiceName = l.cfg.App.Name
	}
	if l.cfg.Observability.OtelExporterOTLPEndpoint == "" {
		// leave empty: Validate will catch when Observability.Enabled is true
	}

	return &l.cfg, nil
}
