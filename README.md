# chat-grpc

Servidor de chat em Go com gRPC bidirecional, salas em memória e suporte opcional a OpenTelemetry (traces e métricas) para ambientes locais ou baseados em containers.


Este repositório foi pensado como um laboratório prático para experimentar e validar padrões importantes em sistemas distribuídos baseados em gRPC. A ideia não é apenas fornecer um servidor de chat funcional, mas oferecer uma base enxuta e repetível para testar:

- Conectividade e streaming bidirecional: como clientes mantêm streams long‑lived com o servidor e como o servidor distribui eventos (join/leave/message) para participantes.  
- Concorrência e coerência de sala: comportamento sob múltiplos clientes concorrentes (ordem de mensagens, perda de mensagens, limpeza de sessões).  
- Resiliência e tratamento de erros: como o sistema reage a cancelamentos, quedas de cliente e reconexões.  
- Observability / Telemetria: como instrumentar o serviço com OpenTelemetry (traces + métricas) e validar que spans e métricas chegam ao collector (OTLP gRPC/HTTP).  
- Integração e deploy em containers: padrões práticos para empacotar o servidor/cliente em imagem Docker e testar redes entre serviços no Docker Compose.

---

## Sumário

- [Visão geral](#visão-geral)
- [Estrutura do repositório](#estrutura-do-repositório)
- [Pré‑requisitos](#pre-requisitos)
- [Execução local](#execução-local)
- [Execução com Docker Compose](#execução-com-docker-compose)
- [Cliente CLI](#cliente-cli)
- [Geração de protos](#geração-de-protos)
- [Testes e qualidade](#testes-e-qualidade)
- [Comandos úteis](#comandos-úteis)

---

## Visão geral

- Serviço gRPC com streaming bidirecional trocando envelopes (`ClientEnvelope` ↔ `ServerEvent`).
- Núcleo de domínio em memória que gerencia salas, sessões e broadcast sem dependências externas.
- Cliente CLI interativo para depuração e demonstrações rápidas.

---

## Estrutura do repositório

- `cmd/chat-grpc` — entrypoint do servidor gRPC.
- `cmd/chat-client` — cliente CLI interativo.
- `api/proto` — definição `chat.proto` e o código gerado em `api/proto/chatv1`.
- `internal/chat` — domínio, salas, sessão, orchestrators e adapters.
- `internal/platform` — config, observability, logger adapter, servidor gRPC.
- `infra/docker` — Dockerfile + ambiente Compose.
- `infra/observability/otel` — configuração do Collector utilizado no Compose.
- `tests/integration` — suíte de integração.
- `makefiles` — alvos especializados usados pelo Makefile principal.

---

## Pre‑requisitos

- Go 1.22+
- Docker 24+ e Docker Compose Plugin (opcional)
- `protoc` instalado para regenerar protos (opcional)

---

## Execução local

1. (Opcional) Garanta caches locais para builds determinísticos:

```bash
export GOCACHE=$PWD/.gocache
export GOMODCACHE=$PWD/.gomodcache
```

2. Suba o servidor diretamente:

```bash
go run ./cmd/chat-grpc
```

3. Abra quantos clientes desejar (cada um em um terminal):

```bash
go run ./cmd/chat-client
```

---

## Execução com Docker Compose

Existem três alvos principais no `Makefile` para o ambiente de desenvolvimento:

- `make build-dev` — constrói a imagem Docker atualizada.
- `make dev` —  sobe o ambiente completo (server + otel collector).
- `make server` — sobe apenas o serviço `chat-grpc`.
- `make client` — abre o cliente CLI dentro do Compose.

Fluxos recomendados:

- Primeira vez / ambiente completo:

```bash
make dev
```

- Iteração rápida se a imagem já existe:

```bash
make server
```

- Forçar rebuild da imagem e subir o servidor:

```bash
make build-dev
make server
```
---

## Cliente CLI

Modo local:

```bash
go run ./cmd/chat-client
```

Modo via compose:

```bash
make client
```
---

## Geração de protos

Se `api/proto/chat.proto` for modificado, regenere os artefatos:

```bash
make proto
```
---

## Testes e qualidade

### Testes unitários

```bash
make test
```

### Testes de integração

Os testes de integração estão em `tests/integration`. Execute com:

```bash
make test-integration
```

---

## Comandos úteis

- `make build` — compila binário do servidor local
- `make run` — executa o servidor local
- `make build-dev` — constrói imagem dev
- `make dev` — sobe ambiente completo (build + up)
- `make server` — sobe apenas o serviço do servidor
- `make client` — abre cliente dentro do compose
- `make test-integration` — roda testes de integração (tag integration)
- `make lint` — roda linters/format
- `make proto` — regenera código a partir do proto

---
