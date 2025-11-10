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
	@grep -q '^CHAT_GRPC_PORT=' .env >/dev/null 2>&1 || echo 'CHAT_GRPC_PORT=50051' >> .env

.PHONY: proto
proto:
	protoc \
		--proto_path=api/proto \
		--go_out=api/proto/chatv1 --go_opt=paths=source_relative \
		--go-grpc_out=api/proto/chatv1 --go-grpc_opt=paths=source_relative \
		chat.proto

# ============================================================
#                DOCKER ENVIRONMENT TARGETS
# ============================================================

APPLICATION_NAME := chat-grpc
DOCKERFILE := infra/docker/Dockerfile
COMPOSE_FILE_DEV := infra/docker/environment/dev/docker-compose-dev.yaml
ENV_FILE := .env

# Build dev image
.PHONY: build-dev
build-dev: clean-dev env
	@echo "\033[1;36m[BUILD-DEV]\033[0m Building DEV image..."
	docker build -f $(DOCKERFILE) -t $(APPLICATION_NAME):dev .

# Start dev environment (builds images as needed)
.PHONY: dev-up
dev-up: dev-down env
	@echo "\033[1;36m[DEV-UP]\033[0m Starting DEV environment..."
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE_DEV) up --build -d

# Stop dev environment
.PHONY: dev-down
dev-down:
	@echo "\033[1;36m[DEV-DOWN]\033[0m Stopping DEV environment..."
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE_DEV) down -v

# Combined target: build image and start dev
.PHONY: dev
dev: build-dev dev-up

# --- Server targets (start/stop only the chat-grpc service) ---
.PHONY: server
server: env
	@echo "\033[1;36m[SERVER]\033[0m Starting chat-grpc service (via docker-compose)..."
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE_DEV) up --build -d chat-grpc

.PHONY: server-down
server-down:
	@echo "\033[1;36m[SERVER-DOWN]\033[0m Stopping chat-grpc service..."
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE_DEV) stop chat-grpc || true
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE_DEV) rm -f chat-grpc || true

# Clean dev artifacts
.PHONY: clean-dev
clean-dev:
	@echo "\033[1;33m[CLEAN-DEV]\033[0m Cleaning DEV containers, volumes, images..."
	@docker ps -a --filter "name=dev" -q | xargs -r docker rm -f || true
	@docker volume ls --filter "name=dev" -q | xargs -r docker volume rm || true
	@docker images --filter "reference=$(APPLICATION_NAME):dev" -q | xargs -r docker rmi -f || true

# --- Client target (docker via compose) ---
.PHONY: client
client: env
	@echo "\033[1;36m[CLIENT]\033[0m Starting a new chat client (compose run)..."
	@docker image inspect dev-chat-grpc >/dev/null 2>&1 || (echo "image not found, building..." && $(MAKE) build-dev)
	@docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE_DEV) run --rm chat-client

# --- Integration tests (gated by build tag) ---
.PHONY: test-integration
test-integration: env
	@echo "\033[1;36m[TEST-INTEGRATION]\033[0m Running integration tests (tag=integration)..."
	@GOCACHE=$(PWD)/.gocache GOMODCACHE=$(PWD)/.gomodcache \
		go test -tags=integration ./tests/integration -v -timeout 60s

# Remove many docker artifacts (use with care)
.PHONY: docker-clean-all
docker-clean-all:
	@echo "\033[1;33m[CLEAN-ALL]\033[0m Removing ALL containers, volumes, images..."
	@docker ps -a -q | xargs -r docker rm -f || true
	@docker volume ls -q | xargs -r docker volume rm || true
	@docker images -a -q | xargs -r docker rmi -f || true
