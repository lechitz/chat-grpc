package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/lechitz/chat-grpc/internal/chat/core/domain"
	"github.com/stretchr/testify/require"
)

type fakeClock struct {
	t time.Time
}

func (f fakeClock) Now() time.Time {
	return f.t
}

func TestJoinCreatesSession(t *testing.T) {
	clk := fakeClock{t: time.Date(2024, 7, 10, 15, 0, 0, 0, time.UTC)}
	svc := NewService(WithClock(clk))

	session, events, err := svc.Join(context.Background(), domain.JoinRequest{
		UserID:      "alice",
		RoomID:      "room-1",
		DisplayName: "",
	})
	require.NoError(t, err)
	require.NotNil(t, events)

	require.Equal(t, domain.Session{
		UserID:      "alice",
		DisplayName: "alice",
		RoomID:      "room-1",
		JoinedAt:    clk.t,
	}, session)
}

func TestJoinDuplicateFails(t *testing.T) {
	svc := NewService()

	_, _, err := svc.Join(context.Background(), domain.JoinRequest{
		UserID: "alice",
		RoomID: "room-1",
	})
	require.NoError(t, err)

	_, _, err = svc.Join(context.Background(), domain.JoinRequest{
		UserID: "alice",
		RoomID: "room-1",
	})
	require.ErrorIs(t, err, ErrAlreadyJoined)
}

func TestBroadcastDeliversToParticipants(t *testing.T) {
	clk := fakeClock{t: time.Date(2024, 7, 10, 15, 0, 0, 0, time.UTC)}
	svc := NewService(WithClock(clk), WithBufferSize(4))

	_, chAlice, err := svc.Join(context.Background(), domain.JoinRequest{
		UserID: "alice",
		RoomID: "room-1",
	})
	require.NoError(t, err)
	require.NotNil(t, chAlice)

	_, chBob, err := svc.Join(context.Background(), domain.JoinRequest{
		UserID: "bob",
		RoomID: "room-1",
	})
	require.NoError(t, err)
	require.NotNil(t, chBob)

	// Drain join notification sent to Alice about Bob joining.
	expectEvent(t, chAlice, domain.EventUserJoined)

	err = svc.Broadcast(context.Background(), domain.Message{
		UserID:  "alice",
		RoomID:  "room-1",
		Content: "hello world",
		SentAt:  clk.t,
	})
	require.NoError(t, err)

	evAlice := expectEvent(t, chAlice, domain.EventMessage)
	require.Equal(t, "hello world", evAlice.Content)
	require.Equal(t, "alice", evAlice.UserID)

	evBob := expectEvent(t, chBob, domain.EventMessage)
	require.Equal(t, "hello world", evBob.Content)
	require.Equal(t, "alice", evBob.UserID)
}

func TestLeaveRemovesUserAndNotifiesOthers(t *testing.T) {
	clk := fakeClock{t: time.Date(2024, 7, 10, 15, 0, 0, 0, time.UTC)}
	svc := NewService(WithClock(clk))

	_, chAlice, err := svc.Join(context.Background(), domain.JoinRequest{
		UserID: "alice",
		RoomID: "room-1",
	})
	require.NoError(t, err)

	_, chBob, err := svc.Join(context.Background(), domain.JoinRequest{
		UserID: "bob",
		RoomID: "room-1",
	})
	require.NoError(t, err)

	expectEvent(t, chAlice, domain.EventUserJoined)

	err = svc.Leave(context.Background(), "room-1", "alice")
	require.NoError(t, err)

	evBob := expectEvent(t, chBob, domain.EventUserLeft)
	require.Equal(t, "alice", evBob.UserID)

	_, ok := <-chAlice
	require.False(t, ok, "alice channel should be closed")
}

func expectEvent(t *testing.T, ch <-chan domain.Event, eventType domain.EventType) domain.Event {
	t.Helper()
	select {
	case ev, ok := <-ch:
		require.True(t, ok, "channel closed unexpectedly")
		require.Equal(t, eventType, ev.Type)
		return ev
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("timed out waiting for event type %v", eventType)
	}
	return domain.Event{}
}
