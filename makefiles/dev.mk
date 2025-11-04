.PHONY: run
run:
	$(GO_ENV) go run ./cmd/chat-grpc

.PHONY: build
build:
	$(GO_ENV) go build ./cmd/chat-grpc

.PHONY: tidy
tidy:
	$(GO_ENV) go mod tidy

.PHONY: env
env:
	@cp -n .env.example .env 2>/dev/null || true

.PHONY: proto
proto:
	protoc \
		--proto_path=api/proto \
		--go_out=api/proto/chatv1 --go_opt=paths=source_relative \
		--go-grpc_out=api/proto/chatv1 --go-grpc_opt=paths=source_relative \
		chat.proto
