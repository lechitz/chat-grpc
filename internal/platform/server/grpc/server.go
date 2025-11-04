// Package grpc provides helpers to build and run the gRPC server.
package grpc

import (
	"context"
	"fmt"
	"net"

	chatv1 "github.com/lechitz/chat-grpc/api/proto/chatv1"
	grpcadapter "github.com/lechitz/chat-grpc/internal/chat/adapter/primary/grpc"
	"github.com/lechitz/chat-grpc/internal/platform/bootstrap"
	"github.com/lechitz/chat-grpc/internal/platform/config"
	"github.com/lechitz/chat-grpc/internal/platform/logger"
	mrt "github.com/lechitz/chat-grpc/internal/platform/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Compose builds the gRPC server and listener.
func Compose(cfg *config.Config, deps *bootstrap.AppDependencies, log logger.Logger) (*grpc.Server, net.Listener, error) {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(cfg.Server.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.Server.MaxSendMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{}),
	}

	server := grpc.NewServer(opts...)

	chatv1.RegisterChatServiceServer(server, grpcadapter.NewServer(deps.ChatService, deps.Logger))

	listener, err := net.Listen("tcp", cfg.Server.Addr())
	if err != nil {
		return nil, nil, fmt.Errorf(errFmtListenTCP, cfg.Server.Addr(), err)
	}

	log.Infow(logMsgServerReady, logFieldAddr, cfg.Server.Addr())
	return server, listener, nil
}

// Run wires the gRPC server into the runtime group.
func Run(ctx context.Context, srv *grpc.Server, lis net.Listener, log logger.Logger) error {
	var group mrt.Group

	group.Add(
		func() error {
			log.Infow(logMsgServerStarting, logFieldAddr, lis.Addr().String())
			return srv.Serve(lis)
		},
		func(_ error) {
			log.Infow(logMsgServerStopping)
			srv.GracefulStop()
		},
	)

	return group.Run(ctx)
}
