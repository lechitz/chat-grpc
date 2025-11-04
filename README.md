# chat-grpc

Implementação em Go de um servidor de chat via gRPC com salas em memória. O objetivo é oferecer um ponto de partida simples, porém estruturado, para aplicações de mensageria que exigem streaming bidirecional.

## Visão geral

- **Transporte**: serviço gRPC (`chat.v1.ChatService`) com um único método bidirecional (`Channel`).
- **Core de domínio**: gerenciador de sessões in-memory em `internal/chat/core/usecase`.
- **Adapter primário**: `internal/chat/adapter/primary/grpc` traduz envelopes protobuf para operações do domínio.
- **Bootstrap**: `cmd/chat-grpc/main.go` integra configuração, logger (`zap`) e servidor.
- **Protobuf**: contrato definido em `api/proto/chat.proto` com bindings Go em `api/proto/chatv1`.

```
.
├── api/proto              # contrato público gRPC
├── cmd/chat-client        # cliente CLI interativo
├── cmd/chat-grpc          # entrypoint do servidor
├── internal/chat          # domínio + adapters
└── internal/platform      # config, logger, runtime, servidor
```

## Pré-requisitos

| Ferramenta | Versão recomendada | Observações |
|------------|--------------------|-------------|
| Go         | 1.22+              | Necessário para build, testes e cliente CLI |
| protoc + plugins Go | opcional | Apenas se desejar regenerar os arquivos em `api/proto/chatv1` |

## Instalação

```bash
git clone https://github.com/lechitz/chat-grpc.git
cd chat-grpc
export GOCACHE=$PWD/.gocache
export GOMODCACHE=$PWD/.gomodcache
go mod download
```

## Configuração

```bash
make env # copia .env.example para .env
```

Edite o arquivo `.env` com os valores adequados antes de levantar o servidor ou o cliente.

## Execução rápida

| Terminal | Comando | Descrição |
|----------|---------|-----------|
| #1 | `make run` | Inicia o servidor gRPC |
| #2 | `make client` | Abre o cliente interativo; informe o nome e a sala (Enter mantém `general`) |
| #3… | `make client` | Inicie quantos participantes quiser; cada terminal mantém um stream ativo |

Dentro do cliente, digite mensagens e pressione Enter. Para encerrar, envie `!quit`. Caso o servidor esteja em outra máquina ou porta, defina `CHAT_GRPC_HOST` e `CHAT_GRPC_PORT` antes de executar `make run`/`make client`.

## Testes

```bash
export GOCACHE=${GOCACHE:-$PWD/.gocache}
export GOMODCACHE=${GOMODCACHE:-$PWD/.gomodcache}
go test ./...
```

## Teste alternativo com `grpcurl`

```bash
printf '%s\n%s\n' \
  '{"join":{"userId":"alice","room":"general","displayName":"Alice"}}' \
  '{"chat":{"userId":"alice","room":"general","content":"Olá pessoal!"}}' \
| grpcurl -plaintext -import-path api/proto -proto chat.proto -d @ \
  127.0.0.1:50051 chat.v1.ChatService.Channel
```

## Regenerar protobuf (opcional)

```bash
protoc \
  --go_out=api/proto/chatv1 --go_opt=paths=source_relative \
  --go-grpc_out=api/proto/chatv1 --go-grpc_opt=paths=source_relative \
  api/proto/chat.proto
```
