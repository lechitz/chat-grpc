package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/lechitz/chat-grpc/internal/chat/core/domain"
)

// Clock abstracts time generation to ease testing.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now().UTC() }

type room struct {
	sessions    map[string]domain.Session
	subscribers map[string]chan domain.Event
}

// Service orchestrates in-memory chat rooms.
type Service struct {
	mu      sync.RWMutex
	rooms   map[string]*room
	clock   Clock
	bufSize int
}

const defaultBufferSize = 32

// Option customises the service behaviour.
type Option func(*Service)

// WithClock overrides the clock used by the service.
func WithClock(clock Clock) Option {
	return func(s *Service) {
		if clock != nil {
			s.clock = clock
		}
	}
}

// WithBufferSize overrides the per-subscriber buffer size.
func WithBufferSize(size int) Option {
	return func(s *Service) {
		if size > 0 {
			s.bufSize = size
		}
	}
}

// NewService creates a new in-memory chat service instance.
func NewService(opts ...Option) *Service {
	svc := &Service{
		rooms:   make(map[string]*room),
		clock:   realClock{},
		bufSize: defaultBufferSize,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Join registers a user in the requested room and returns a session plus the event stream.
func (s *Service) Join(_ context.Context, req domain.JoinRequest) (domain.Session, <-chan domain.Event, error) {
	if req.RoomID == "" || req.UserID == "" {
		return domain.Session{}, nil, ErrEmptyFields
	}

	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.UserID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	rm := s.ensureRoom(req.RoomID)
	if _, exists := rm.sessions[req.UserID]; exists {
		return domain.Session{}, nil, ErrAlreadyJoined
	}

	session := domain.Session{
		UserID:      req.UserID,
		DisplayName: displayName,
		RoomID:      req.RoomID,
		JoinedAt:    s.clock.Now(),
	}
	eventCh := make(chan domain.Event, s.bufSize)

	rm.sessions[req.UserID] = session
	rm.subscribers[req.UserID] = eventCh

	s.enqueueLocked(req.RoomID, domain.Event{
		Type:        domain.EventUserJoined,
		UserID:      session.UserID,
		DisplayName: session.DisplayName,
		RoomID:      session.RoomID,
		Timestamp:   session.JoinedAt,
	}, req.UserID)

	return session, eventCh, nil
}

// Leave unregisters the user and notifies remaining participants.
func (s *Service) Leave(_ context.Context, roomID, userID string) error {
	if roomID == "" || userID == "" {
		return ErrEmptyFields
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	rm, ok := s.rooms[roomID]
	if !ok {
		return ErrRoomNotFound
	}

	session, ok := rm.sessions[userID]
	if !ok {
		return ErrUserNotInRoom
	}

	ch := rm.subscribers[userID]
	delete(rm.sessions, userID)
	delete(rm.subscribers, userID)
	close(ch)

	s.enqueueLocked(roomID, domain.Event{
		Type:        domain.EventUserLeft,
		UserID:      session.UserID,
		DisplayName: session.DisplayName,
		RoomID:      session.RoomID,
		Timestamp:   s.clock.Now(),
	}, userID)

	if len(rm.sessions) == 0 {
		delete(s.rooms, roomID)
	}

	return nil
}

// Broadcast delivers a message to all participants in the room.
func (s *Service) Broadcast(_ context.Context, msg domain.Message) error {
	if msg.RoomID == "" || msg.UserID == "" {
		return ErrEmptyFields
	}
	if msg.Content == "" {
		return ErrEmptyMessage
	}

	s.mu.RLock()
	rm, ok := s.rooms[msg.RoomID]
	if !ok {
		s.mu.RUnlock()
		return ErrRoomNotFound
	}

	session, ok := rm.sessions[msg.UserID]
	if !ok {
		s.mu.RUnlock()
		return ErrUserNotInRoom
	}

	event := domain.Event{
		Type:        domain.EventMessage,
		UserID:      session.UserID,
		DisplayName: session.DisplayName,
		RoomID:      session.RoomID,
		Content:     msg.Content,
		Timestamp:   s.clock.Now(),
	}

	channels := make([]chan domain.Event, 0, len(rm.subscribers))
	for _, ch := range rm.subscribers {
		channels = append(channels, ch)
	}
	s.mu.RUnlock()

	for _, ch := range channels {
		select {
		case ch <- event:
		default:
		}
	}

	return nil
}

func (s *Service) ensureRoom(roomID string) *room {
	rm, ok := s.rooms[roomID]
	if !ok {
		rm = &room{
			sessions:    make(map[string]domain.Session),
			subscribers: make(map[string]chan domain.Event),
		}
		s.rooms[roomID] = rm
	}
	return rm
}

// enqueueLocked broadcasts an event to a room while holding the global lock.
func (s *Service) enqueueLocked(roomID string, event domain.Event, excludeUser string) {
	rm := s.rooms[roomID]
	if rm == nil {
		return
	}
	for uid, ch := range rm.subscribers {
		if uid == excludeUser {
			continue
		}
		select {
		case ch <- event:
		default:
		}
	}
}
