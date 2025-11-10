package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lechitz/chat-grpc/internal/adapter/secondary/contextlogger"
	"github.com/lechitz/chat-grpc/internal/platform/bootstrap"
	"github.com/lechitz/chat-grpc/internal/platform/config"
	"github.com/lechitz/chat-grpc/internal/platform/observability/metric"
	"github.com/lechitz/chat-grpc/internal/platform/observability/tracer"
	"github.com/lechitz/chat-grpc/internal/platform/ports/logger"
	"github.com/lechitz/chat-grpc/internal/platform/server"
	"github.com/lechitz/chat-grpc/internal/shared/commonkeys"
)

func main() {
	os.Exit(run())
}

func run() int {
	logs, cleanupLogger := contextlogger.New()
	defer cleanupLogger()

	cfg, err := loadConfig(logs)
	if err != nil {
		logs.Errorw(ErrLoadConfig, commonkeys.Error, err.Error())
		return 2
	}

	tracerCleanup, err := tracer.Init(cfg, logs)
	if err != nil {
		logs.Errorw(logMsgInitObservability, logFieldError, err)
		return 2
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ServerGRPC.ShutdownGrace)
		defer cancel()
		tracerCleanup(ctx)
	}()

	metricCleanup, err := metric.Init(cfg, logs)
	if err != nil {
		logs.Errorw(logMsgInitObservability, logFieldError, err)
		return 2
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ServerGRPC.ShutdownGrace)
		defer cancel()
		metricCleanup(ctx)
	}()

	logs.Infow(logMsgConfigLoaded,
		logFieldApp, cfg.App.Name,
		logFieldEnv, cfg.App.Environment,
		logFieldAddr, cfg.ServerGRPC.Addr(),
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	deps, cleanup, err := bootstrap.Initialize(ctx, cfg, logs)
	if err != nil {
		logs.Errorw(logMsgInitializeDeps, logFieldError, err)
		return 3
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ServerGRPC.ShutdownGrace)
		defer cancel()
		cleanup(shutdownCtx)
	}()

	if err := server.RunAll(ctx, cfg, deps, logs); err != nil {
		logs.Errorw(logMsgServerFailed, logFieldError, err)
		return 1
	}

	return 0
}

// loadConfig loads the application configuration.
func loadConfig(logs logger.ContextLogger) (*config.Config, error) {
	cfg, err := config.New().Load(logs)
	if err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	logs.Infow(
		MsgConfigLoaded,
		commonkeys.APIName, cfg.App.Name,
		commonkeys.AppEnv, cfg.App.Environment,
		commonkeys.AppVersion, cfg.General.Version,
	)
	return cfg, nil
}
