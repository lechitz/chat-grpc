package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	chatv1 "github.com/lechitz/chat-grpc/api/proto/chatv1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	displayName, err := prompt(reader, promptDisplayName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao ler nome: %v\n", err)
		return 1
	}
	displayName = sanitizeInput(displayName)
	if displayName == "" {
		displayName = generateFallbackName()
	}

	room, err := prompt(reader, promptRoom)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao ler sala: %v\n", err)
		return 1
	}
	room = sanitizeInput(room)
	if room == "" {
		room = defaultRoom
	}

	userID := buildUserID(displayName)
	target := fmt.Sprintf("%s:%s", getenv(envHostKey, defaultHost), getenv(envPortKey, defaultPort))

	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		fmt.Fprintf(os.Stderr, "falha ao conectar a %s: %v\n", target, err)
		return 2
	}
	defer conn.Close()

	client := chatv1.NewChatServiceClient(conn)
	stream, err := client.Channel(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "falha ao abrir stream: %v\n", err)
		return 2
	}

	if err := stream.Send(&chatv1.ClientEnvelope{
		Message: &chatv1.ClientEnvelope_Join{
			Join: &chatv1.JoinRequest{
				UserId:      userID,
				DisplayName: displayName,
				Room:        room,
			},
		},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "falha ao enviar join: %v\n", err)
		return 2
	}

	firstEvent, err := stream.Recv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "falha ao receber ack: %v\n", err)
		return 2
	}
	if ack := firstEvent.GetJoined(); ack != nil {
		fmt.Printf(messageConnected+"\n", ack.GetRoom(), displayName)
		fmt.Println(messagePromptCommands)
	} else {
		fmt.Println(messageInvalidJoinAck)
		return 2
	}

	var wg sync.WaitGroup
	recvDone := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(recvDone)
		for {
			event, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
					fmt.Println("\n" + messageServerClosed)
				} else {
					fmt.Printf("\n"+messageReceiveError, err)
				}
				return
			}
			renderEvent(event)
			fmt.Print(promptInput)
		}
	}()

	for {
		fmt.Print(promptInput)
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				line = strings.TrimSpace(line)
				if line == "" {
					line = commandQuit
				}
			} else {
				fmt.Printf(messageSendError, err)
				continue
			}
		}

		line = sanitizeInput(line)
		if line == "" {
			continue
		}
		if line == commandQuit {
			fmt.Println(messageLeaving)
			if err := stream.Send(&chatv1.ClientEnvelope{
				Message: &chatv1.ClientEnvelope_Leave{
					Leave: &chatv1.LeaveRequest{
						UserId: userID,
						Room:   room,
					},
				},
			}); err != nil {
				fmt.Printf(messageLeaveError, err)
			}
			_ = stream.CloseSend()
			<-recvDone
			wg.Wait()
			fmt.Println(messageDisconnected)
			return 0
		}

		if err := stream.Send(&chatv1.ClientEnvelope{
			Message: &chatv1.ClientEnvelope_Chat{
				Chat: &chatv1.ChatPayload{
					UserId:       userID,
					Room:         room,
					Content:      line,
					TimestampUtc: time.Now().UTC().UnixMilli(),
				},
			},
		}); err != nil {
			fmt.Printf(messageSendError, err)
		}
	}
}

func prompt(reader *bufio.Reader, question string) (string, error) {
	fmt.Print(question)
	text, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func sanitizeInput(value string) string {
	value = strings.TrimSpace(value)
	return strings.ToValidUTF8(value, "")
}

func buildUserID(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return generateFallbackName()
	}
	var b strings.Builder
	lastHyphen := false
	for _, r := range name {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastHyphen = false
		case unicode.IsSpace(r) || r == '-' || r == '_':
			if !lastHyphen {
				b.WriteRune('-')
				lastHyphen = true
			}
		}
	}
	id := strings.Trim(b.String(), "-")
	if id == "" {
		return generateFallbackName()
	}
	return id
}

func generateFallbackName() string {
	return fmt.Sprintf("user-%d", time.Now().UnixNano())
}

func getenv(key, fallback string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return fallback
}

func renderEvent(event *chatv1.ServerEvent) {
	switch payload := event.GetEvent().(type) {
	case *chatv1.ServerEvent_Broadcast:
		if payload.Broadcast == nil {
			return
		}
		tsVal := payload.Broadcast.GetTimestampUtc()
		timestamp := time.Now().UTC()
		if tsVal != 0 {
			timestamp = time.UnixMilli(tsVal)
		}
		fmt.Printf(messageIncomingChat+"\n", timestamp.Format(timeDisplayFormat), payload.Broadcast.GetUserId(), strings.ToValidUTF8(payload.Broadcast.GetContent(), ""))
	case *chatv1.ServerEvent_Notice:
		renderNotice(payload.Notice)
	default:
		fmt.Println(messageUnknownEvent)
	}
}

func renderNotice(notice *chatv1.ServerNotice) {
	if notice == nil {
		return
	}
	switch notice.GetType() {
	case chatv1.ServerNotice_TYPE_USER_JOINED:
		fmt.Printf(messageNoticeUserJoined+"\n", displayNameFallback(notice.GetUserId()))
	case chatv1.ServerNotice_TYPE_USER_LEFT:
		fmt.Printf(messageNoticeUserLeft+"\n", displayNameFallback(notice.GetUserId()))
	case chatv1.ServerNotice_TYPE_GENERIC:
		fmt.Printf(messageNoticeGeneric+"\n", notice.GetMessage())
	case chatv1.ServerNotice_TYPE_ERROR:
		fmt.Printf(messageSystemError+"\n", notice.GetMessage())
	default:
		fmt.Printf(messageNoticeGeneric+"\n", notice.GetMessage())
	}
}

func displayNameFallback(userID string) string {
	if userID == "" {
		return "usuário"
	}
	return userID
}
