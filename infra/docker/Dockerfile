FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the chat-grpc server and chat-client CLI
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chat-grpc ./cmd/chat-grpc
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chat-client ./cmd/chat-client

FROM alpine:3.19.1
RUN apk add --no-cache ca-certificates
WORKDIR /app

# Copy the compiled binaries
COPY --from=builder /app/chat-grpc .
COPY --from=builder /app/chat-client .

# Default to running the server; override when running the client container
ENTRYPOINT ["./chat-grpc"]
