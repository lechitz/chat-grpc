package config

import "time"

const (
	envFileName = ".env"

	envAppNameKey            = "CHAT_GRPC_APP_NAME"
	envEnvironmentKey        = "CHAT_GRPC_ENV"
	envHostKey               = "CHAT_GRPC_HOST"
	envPortKey               = "CHAT_GRPC_PORT"
	envShutdownGraceKey      = "CHAT_GRPC_SHUTDOWN_GRACE"
	envMaxRecvSizeKey        = "CHAT_GRPC_MAX_RECV_MSG_SIZE"
	envMaxSendSizeKey        = "CHAT_GRPC_MAX_SEND_MSG_SIZE"
	envOtelEnabledKey        = "CHAT_GRPC_OTEL_ENABLED"
	envOtelEndpointKey       = "CHAT_GRPC_OTEL_EXPORTER_ENDPOINT"
	envOtelInsecureKey       = "CHAT_GRPC_OTEL_EXPORTER_INSECURE"
	envOtelTimeoutKey        = "CHAT_GRPC_OTEL_EXPORTER_TIMEOUT"
	envOtelCompressionKey    = "CHAT_GRPC_OTEL_EXPORTER_COMPRESSION"
	envOtelHeadersKey        = "CHAT_GRPC_OTEL_EXPORTER_HEADERS"
	envOtelServiceNameKey    = "CHAT_GRPC_OTEL_SERVICE_NAME"
	envOtelServiceVersionKey = "CHAT_GRPC_OTEL_SERVICE_VERSION"

	defaultAppName            = "chat-grpc"
	defaultEnvironment        = "development"
	defaultHost               = "127.0.0.1"
	defaultPort               = "50051"
	defaultShutdownGrace      = 5 * time.Second
	defaultMaxRecvMsgSize     = 4 << 20 // 4 MiB
	defaultMaxSendMsgSize     = 4 << 20 // 4 MiB
	defaultOtelEnabled        = false
	defaultOtelInsecure       = true
	defaultOtelTimeout        = "5s"
	defaultOtelCompression    = "none"
	defaultOtelServiceVersion = "0.1.0"
)
