// Package domain holds the business entities used by the chat use cases.
package domain

import "time"

// JoinRequest represents a user attempting to join a chat room.
type JoinRequest struct {
	UserID      string
	DisplayName string
	RoomID      string
}

// Session describes an active connection inside a room.
type Session struct {
	UserID      string
	DisplayName string
	RoomID      string
	JoinedAt    time.Time
}

// Message is the canonical event broadcast to room participants.
type Message struct {
	UserID      string
	DisplayName string
	RoomID      string
	Content     string
	SentAt      time.Time
}

// EventType categorizes outbound events delivered to participants.
type EventType int

const (
	// EventMessage indicates a chat message broadcast.
	EventMessage EventType = iota
	// EventUserJoined indicates someone joined the room.
	EventUserJoined
	// EventUserLeft indicates someone left the room.
	EventUserLeft
	// EventSystem carries generic notices, typically errors.
	EventSystem
)

// Event represents a server-side notification pushed to clients.
type Event struct {
	Type        EventType
	UserID      string
	DisplayName string
	RoomID      string
	Content     string
	Timestamp   time.Time
}
