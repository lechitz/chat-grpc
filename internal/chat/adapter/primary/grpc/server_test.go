package grpcadapter_test

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	chatv1 "github.com/lechitz/chat-grpc/api/proto/chatv1"
	grpcadapter "github.com/lechitz/chat-grpc/internal/chat/adapter/primary/grpc"
	"github.com/lechitz/chat-grpc/internal/chat/core/usecase"
	"github.com/lechitz/chat-grpc/internal/platform/logger"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type noopLogger struct{}

func (noopLogger) Debugw(string, ...interface{}) {}
func (noopLogger) Infow(string, ...interface{})  {}
func (noopLogger) Warnw(string, ...interface{})  {}
func (noopLogger) Errorw(string, ...interface{}) {}
func (noopLogger) Sync() error                   { return nil }

var _ logger.Logger = (*noopLogger)(nil)

func TestChannel_BasicFlow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	t.Cleanup(srv.Stop)

	app := usecase.NewService()
	chatv1.RegisterChatServiceServer(srv, grpcadapter.NewServer(app, noopLogger{}))

	go func() {
		_ = srv.Serve(lis)
	}()

	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client := chatv1.NewChatServiceClient(conn)
	stream, err := client.Channel(ctx)
	require.NoError(t, err)

	err = stream.Send(&chatv1.ClientEnvelope{
		Message: &chatv1.ClientEnvelope_Join{
			Join: &chatv1.JoinRequest{
				UserId:      "alice",
				Room:        "general",
				DisplayName: "Alice",
			},
		},
	})
	require.NoError(t, err)

	ev, err := stream.Recv()
	require.NoError(t, err)
	ack := ev.GetJoined()
	require.NotNil(t, ack)
	require.Equal(t, "alice", ack.GetUserId())
	require.Equal(t, "general", ack.GetRoom())

	err = stream.Send(&chatv1.ClientEnvelope{
		Message: &chatv1.ClientEnvelope_Chat{
			Chat: &chatv1.ChatPayload{
				UserId:  "alice",
				Room:    "general",
				Content: "olá mundo",
			},
		},
	})
	require.NoError(t, err)

	var broadcast *chatv1.ServerEvent
	assertWithin(t, time.Second, func() bool {
		broadcast, err = stream.Recv()
		if err != nil {
			return false
		}
		return broadcast.GetBroadcast() != nil
	})
	require.NoError(t, err)
	require.NotNil(t, broadcast.GetBroadcast())
	require.Equal(t, "olá mundo", broadcast.GetBroadcast().GetContent())

	err = stream.Send(&chatv1.ClientEnvelope{
		Message: &chatv1.ClientEnvelope_Leave{
			Leave: &chatv1.LeaveRequest{
				UserId: "alice",
				Room:   "general",
			},
		},
	})
	require.NoError(t, err)

	_, err = stream.Recv()
	require.Error(t, err)
	require.True(t, errors.Is(err, io.EOF) || errors.Is(err, context.Canceled))
}

func assertWithin(t *testing.T, timeout time.Duration, fn func() bool) {
	t.Helper()
	expire := time.Now().Add(timeout)
	for time.Now().Before(expire) {
		if fn() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("condition not satisfied within %s", timeout)
}
