// Code generated manually to provide gRPC bindings for chat.proto. DO NOT EDIT.
package chatv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// ChatServiceClient is the client API for ChatService service.
type ChatServiceClient interface {
	Channel(ctx context.Context, opts ...grpc.CallOption) (ChatService_ChannelClient, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

// NewChatServiceClient creates a new client.
func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) Channel(ctx context.Context, opts ...grpc.CallOption) (ChatService_ChannelClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[0], "/chat.v1.ChatService/Channel", opts...)
	if err != nil {
		return nil, err
	}
	return &chatServiceChannelClient{stream}, nil
}

type ChatService_ChannelClient interface {
	Send(*ClientEnvelope) error
	Recv() (*ServerEvent, error)
	grpc.ClientStream
}

type chatServiceChannelClient struct {
	grpc.ClientStream
}

func (x *chatServiceChannelClient) Send(m *ClientEnvelope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServiceChannelClient) Recv() (*ServerEvent, error) {
	m := new(ServerEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServiceServer is the server API for ChatService service.
type ChatServiceServer interface {
	Channel(ChatService_ChannelServer) error
}

// UnimplementedChatServiceServer should be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct{}

func (UnimplementedChatServiceServer) Channel(ChatService_ChannelServer) error {
	return status.Errorf(codes.Unimplemented, "method Channel not implemented")
}

func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// This is not recommended and should be used only by legacy code.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_Channel_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServiceServer).Channel(&chatServiceChannelServer{stream})
}

type ChatService_ChannelServer interface {
	Send(*ServerEvent) error
	Recv() (*ClientEnvelope, error)
	grpc.ServerStream
}

type chatServiceChannelServer struct {
	grpc.ServerStream
}

func (x *chatServiceChannelServer) Send(m *ServerEvent) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServiceChannelServer) Recv() (*ClientEnvelope, error) {
	m := new(ClientEnvelope)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.v1.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Channel",
			Handler:       _ChatService_Channel_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "chat.proto",
}
