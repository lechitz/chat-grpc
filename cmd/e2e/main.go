package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	chatv1 "github.com/lechitz/chat-grpc/api/proto/chatv1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	host := getenv("CHAT_GRPC_HOST", "127.0.0.1")
	port := getenv("CHAT_GRPC_PORT", "50051")
	target := fmt.Sprintf("%s:%s", host, port)

	names := []string{"alice", "bob", "carol"}
	var wg sync.WaitGroup
	wg.Add(len(names))

	for i, name := range names {
		go func(idx int, user string) {
			defer wg.Done()
			runClient(target, user, "general", fmt.Sprintf("hello from %s", user), time.Duration(500*idx)*time.Millisecond)
		}(i, name)
	}

	wg.Wait()
}

func runClient(target, name, room, message string, startDelay time.Duration) {
	ctx := context.Background()
	if startDelay > 0 {
		time.Sleep(startDelay)
	}

	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		fmt.Printf("%s: failed to dial: %v\n", name, err)
		return
	}
	defer conn.Close()

	client := chatv1.NewChatServiceClient(conn)
	stream, err := client.Channel(ctx)
	if err != nil {
		fmt.Printf("%s: failed to open stream: %v\n", name, err)
		return
	}

	// send join
	tjoin := &chatv1.ClientEnvelope{Message: &chatv1.ClientEnvelope_Join{Join: &chatv1.JoinRequest{UserId: name, DisplayName: name, Room: room}}}
	if err := stream.Send(tjoin); err != nil {
		fmt.Printf("%s: send join error: %v\n", name, err)
		return
	}

	// receive joined ack
	ack, err := stream.Recv()
	if err != nil {
		fmt.Printf("%s: recv join ack error: %v\n", name, err)
		return
	}
	if j := ack.GetJoined(); j != nil {
		fmt.Printf("%s: joined room %s\n", name, j.Room)
	} else {
		fmt.Printf("%s: unexpected first event: %v\n", name, ack)
	}

	// start receiver
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			evt, err := stream.Recv()
			if err != nil {
				return
			}
			if b := evt.GetBroadcast(); b != nil {
				fmt.Printf("%s received broadcast: %s: %s\n", name, b.UserId, b.Content)
			} else if n := evt.GetNotice(); n != nil {
				fmt.Printf("%s received notice: %s\n", name, n.Message)
			}
		}
	}()

	// wait small time to allow others to join
	time.Sleep(500 * time.Millisecond)

	// send chat
	if err := stream.Send(&chatv1.ClientEnvelope{Message: &chatv1.ClientEnvelope_Chat{Chat: &chatv1.ChatPayload{UserId: name, Room: room, Content: message, TimestampUtc: time.Now().UTC().UnixMilli()}}}); err != nil {
		fmt.Printf("%s: send chat error: %v\n", name, err)
	}

	// wait to receive broadcasts
	time.Sleep(800 * time.Millisecond)

	// leave
	if err := stream.Send(&chatv1.ClientEnvelope{Message: &chatv1.ClientEnvelope_Leave{Leave: &chatv1.LeaveRequest{UserId: name, Room: room}}}); err != nil {
		fmt.Printf("%s: send leave error: %v\n", name, err)
	}

	// close send and wait receiver
	_ = stream.CloseSend()
	<-done
	fmt.Printf("%s: finished\n", name)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
