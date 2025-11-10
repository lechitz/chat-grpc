//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	chatv1 "github.com/lechitz/chat-grpc/api/proto/chatv1"
	"github.com/lechitz/chat-grpc/internal/platform/bootstrap"
	"github.com/lechitz/chat-grpc/internal/platform/config"
	"github.com/lechitz/chat-grpc/internal/platform/logger"
	grpcserver "github.com/lechitz/chat-grpc/internal/platform/server"
	grpcgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// helper to get a free port on localhost
func freeListenPort(t *testing.T) (string, func()) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	addr := l.Addr().(*net.TCPAddr)
	port := addr.Port
	// close the listener and return cleanup
	cleanup := func() { _ = l.Close() }
	return fmt.Sprintf("127.0.0.1:%d", port), cleanup
}

func TestChatFlow(t *testing.T) {
	// load config and override to use localhost dynamic port
	loader := config.New()
	cfg, err := loader.Load(logger.NoopLogger{})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	// choose a free port by listening, then closing
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to obtain free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	cfg.ServerGRPC.Host = "127.0.0.1"
	cfg.ServerGRPC.Port = fmt.Sprintf("%d", port)

	logs := logger.NoopLogger{}
	deps, cleanup, err := bootstrap.Initialize(context.Background(), cfg, logs)
	if err != nil {
		t.Fatalf("bootstrap initialize: %v", err)
	}
	defer cleanup(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run server
	go func() {
		if err := grpcserver.RunAll(ctx, cfg, deps, logs); err != nil {
			// cause test failure if server returns error
			t.Logf("server stopped: %v", err)
		}
	}()

	// wait a bit for server to bind
	tryUntil(t, 3*time.Second, func() bool {
		c, err := grpcgrpc.DialContext(context.Background(), cfg.ServerGRPC.Addr(), grpcgrpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			_ = c.Close()
			return true
		}
		return false
	})

	// create three clients and have them join/send/leave
	users := []string{"alice", "bob", "carol"}
	var wg sync.WaitGroup
	wg.Add(len(users))

	results := make(chan string, 10)

	for i, user := range users {
		go func(u string, delay time.Duration) {
			defer wg.Done()
			conn, err := grpcgrpc.DialContext(context.Background(), cfg.ServerGRPC.Addr(), grpcgrpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				results <- fmt.Sprintf("dial failed: %v", err)
				return
			}
			defer conn.Close()
			client := chatv1.NewChatServiceClient(conn)
			stream, err := client.Channel(context.Background())
			if err != nil {
				results <- fmt.Sprintf("open stream: %v", err)
				return
			}

			// send join
			if err := stream.Send(&chatv1.ClientEnvelope{Message: &chatv1.ClientEnvelope_Join{Join: &chatv1.JoinRequest{UserId: u, DisplayName: u, Room: "general"}}}); err != nil {
				results <- fmt.Sprintf("join send: %v", err)
				return
			}
			// expect joined ack
			if ev, err := stream.Recv(); err != nil {
				results <- fmt.Sprintf("recv ack: %v", err)
				return
			} else if ev.GetJoined() == nil {
				results <- fmt.Sprintf("expected joined ack, got: %v", ev)
				return
			}

			// wait for others then send a message
			time.Sleep(delay)
			if err := stream.Send(&chatv1.ClientEnvelope{Message: &chatv1.ClientEnvelope_Chat{Chat: &chatv1.ChatPayload{UserId: u, Room: "general", Content: "hello from " + u, TimestampUtc: time.Now().UTC().UnixMilli()}}}); err != nil {
				results <- fmt.Sprintf("chat send: %v", err)
				return
			}

			// read a couple of events, then leave
			for k := 0; k < 3; k++ {
				if _, err := stream.Recv(); err != nil {
					break
				}
			}
			if err := stream.Send(&chatv1.ClientEnvelope{Message: &chatv1.ClientEnvelope_Leave{Leave: &chatv1.LeaveRequest{UserId: u, Room: "general"}}}); err != nil {
				results <- fmt.Sprintf("leave send: %v", err)
				return
			}
			// close send
			_ = stream.CloseSend()
			results <- fmt.Sprintf("done %s", u)
		}(user, time.Duration(i)*200*time.Millisecond)
	}

	wg.Wait()
	close(results)

	for r := range results {
		if r != "done alice" && r != "done bob" && r != "done carol" {
			t.Fatalf("unexpected result: %s", r)
		}
	}
}

// helper that retries until timeout
func tryUntil(t *testing.T, timeout time.Duration, f func() bool) {
	t0 := time.Now()
	for {
		if f() {
			return
		}
		if time.Since(t0) > timeout {
			t.Fatalf("condition not met within %v", timeout)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
