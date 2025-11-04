package config

import "time"

const (
	envFileName = ".env"

	envAppNameKey       = "CHAT_GRPC_APP_NAME"
	envEnvironmentKey   = "CHAT_GRPC_ENV"
	envHostKey          = "CHAT_GRPC_HOST"
	envPortKey          = "CHAT_GRPC_PORT"
	envShutdownGraceKey = "CHAT_GRPC_SHUTDOWN_GRACE"
	envMaxRecvSizeKey   = "CHAT_GRPC_MAX_RECV_MSG_SIZE"
	envMaxSendSizeKey   = "CHAT_GRPC_MAX_SEND_MSG_SIZE"

	defaultAppName        = "chat-grpc"
	defaultEnvironment    = "development"
	defaultHost           = "127.0.0.1"
	defaultPort           = "50051"
	defaultShutdownGrace  = 5 * time.Second
	defaultMaxRecvMsgSize = 4 << 20 // 4 MiB
	defaultMaxSendMsgSize = 4 << 20 // 4 MiB
)
