// Package server exposes orchestration helpers for starting platform servers.
package server

import (
	"context"
	"fmt"

	"github.com/lechitz/chat-grpc/internal/platform/bootstrap"
	"github.com/lechitz/chat-grpc/internal/platform/config"
	"github.com/lechitz/chat-grpc/internal/platform/ports/logger"
	grpcserver "github.com/lechitz/chat-grpc/internal/platform/server/grpc"
)

// RunAll currently boots only the gRPC server, keeping the signature open for future transports.
func RunAll(ctx context.Context, cfg *config.Config, deps *bootstrap.AppDependencies, log logger.ContextLogger) error {
	srv, lis, err := grpcserver.Compose(cfg, deps, log)
	if err != nil {
		return fmt.Errorf(errFmtComposeGRPCServer, err)
	}

	return grpcserver.Run(ctx, srv, lis, log)
}
