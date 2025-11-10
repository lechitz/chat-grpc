// Package bootstrap wires infrastructure components and domain services.
package bootstrap

import (
	"context"

	"github.com/lechitz/chat-grpc/internal/chat/core/ports/input"
	"github.com/lechitz/chat-grpc/internal/chat/core/usecase"
	"github.com/lechitz/chat-grpc/internal/platform/config"
	"github.com/lechitz/chat-grpc/internal/platform/ports/logger"
)

// AppDependencies collects the primary ports exposed to adapters.
type AppDependencies struct {
	ChatService input.StreamService
	Logger      logger.ContextLogger
}

// Initialize builds the dependencies required by transports.
func Initialize(ctx context.Context, cfg *config.Config, log logger.ContextLogger) (*AppDependencies, func(context.Context), error) {
	_ = cfg
	_ = ctx

	chatService := usecase.NewService()

	// Placeholder for future resource cleanup.
	cleanup := func(context.Context) {}

	return &AppDependencies{
		ChatService: chatService,
		Logger:      log,
	}, cleanup, nil
}
