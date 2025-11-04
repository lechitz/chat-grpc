package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lechitz/chat-grpc/internal/platform/bootstrap"
	"github.com/lechitz/chat-grpc/internal/platform/config"
	"github.com/lechitz/chat-grpc/internal/platform/logger"
	"github.com/lechitz/chat-grpc/internal/platform/server"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.Load()
	if err != nil {
		// nothing we can do if logging is not yet available
		_, _ = os.Stderr.WriteString(stderrLoadConfigPref + err.Error())
		return 2
	}

	logs, cleanupLogger, err := logger.New(cfg.App.Environment)
	if err != nil {
		_, _ = os.Stderr.WriteString(stderrInitLoggerPref + err.Error())
		return 2
	}
	defer cleanupLogger()

	logs.Infow(logMsgConfigLoaded,
		logFieldApp, cfg.App.Name,
		logFieldEnv, cfg.App.Environment,
		logFieldAddr, cfg.Server.Addr(),
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	deps, cleanup, err := bootstrap.Initialize(ctx, cfg, logs)
	if err != nil {
		logs.Errorw(logMsgInitializeDeps, logFieldError, err)
		return 3
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownGrace)
		defer cancel()
		cleanup(shutdownCtx)
	}()

	if err := server.RunAll(ctx, cfg, deps, logs); err != nil {
		logs.Errorw(logMsgServerFailed, logFieldError, err)
		return 1
	}

	return 0
}
