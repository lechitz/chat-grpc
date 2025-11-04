// Package input defines the primary (driving) ports for the chat domain.
package input

import (
	"context"

	"github.com/lechitz/chat-grpc/internal/chat/core/domain"
)

// StreamService exposes the operations consumed by the gRPC adapter.
type StreamService interface {
	Join(ctx context.Context, req domain.JoinRequest) (domain.Session, <-chan domain.Event, error)
	Leave(ctx context.Context, roomID, userID string) error
	Broadcast(ctx context.Context, msg domain.Message) error
}
