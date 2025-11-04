package usecase

import "errors"

var (
	// ErrEmptyFields indicates the join request lacks mandatory information.
	ErrEmptyFields = errors.New("join request missing required fields")
	// ErrAlreadyJoined indicates the user is already present in the room.
	ErrAlreadyJoined = errors.New("user already joined room")
	// ErrRoomNotFound indicates the room is unknown.
	ErrRoomNotFound = errors.New("room not found")
	// ErrUserNotInRoom indicates an operation was attempted for a missing user.
	ErrUserNotInRoom = errors.New("user not part of room")
	// ErrEmptyMessage indicates the message body is empty.
	ErrEmptyMessage = errors.New("message content is empty")
)
