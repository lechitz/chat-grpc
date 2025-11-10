# Local client helper (kept as client-local to avoid overriding the docker client target)
.PHONY: client-local
client-local:
	$(GO_ENV) go run ./cmd/chat-client
