package grpcadapter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	chatv1 "github.com/lechitz/chat-grpc/api/proto/chatv1"
	"github.com/lechitz/chat-grpc/internal/chat/core/domain"
	"github.com/lechitz/chat-grpc/internal/chat/core/ports/input"
	"github.com/lechitz/chat-grpc/internal/chat/core/usecase"
	"github.com/lechitz/chat-grpc/internal/platform/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the generated gRPC ChatServiceServer.
type Server struct {
	chatv1.UnimplementedChatServiceServer

	chat input.StreamService
	log  logger.Logger
}

// NewServer constructs a gRPC adapter backed by the domain chat service.
func NewServer(chat input.StreamService, log logger.Logger) *Server {
	return &Server{
		chat: chat,
		log:  log,
	}
}

// Channel handles the bidirectional chat stream lifecycle.
func (s *Server) Channel(stream chatv1.ChatService_ChannelServer) error {
	ctx := stream.Context()

	var (
		session      domain.Session
		events       <-chan domain.Event
		hasSession   bool
		sendMu       sync.Mutex
		eventsCancel context.CancelFunc
		eventsWG     sync.WaitGroup
		eventErr     = make(chan error, 1)
	)

	send := func(evt *chatv1.ServerEvent) error {
		if evt == nil {
			return nil
		}
		sendMu.Lock()
		defer sendMu.Unlock()
		if err := stream.Send(evt); err != nil {
			return err
		}
		return nil
	}

	defer func() {
		if eventsCancel != nil {
			eventsCancel()
		}
		eventsWG.Wait()
		if hasSession {
			if err := s.chat.Leave(context.Background(), session.RoomID, session.UserID); err != nil && !errors.Is(err, usecase.ErrUserNotInRoom) {
				s.log.Warnw(logMsgCleanupSessionFailure, logFieldRoom, session.RoomID, logFieldUser, session.UserID, logFieldError, err)
			}
		}
	}()

	for {
		select {
		case err := <-eventErr:
			if err != nil {
				return err
			}
		default:
		}

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return status.Error(codes.Canceled, errMsgClientCanceled)
			}
			return err
		}

		switch msg := req.GetMessage().(type) {
		case *chatv1.ClientEnvelope_Join:
			if hasSession {
				return status.Error(codes.FailedPrecondition, errMsgSessionExists)
			}

			in := msg.Join
			if in == nil {
				return status.Error(codes.InvalidArgument, errMsgJoinPayloadRequired)
			}

			session, events, err = s.chat.Join(ctx, domain.JoinRequest{
				UserID:      in.GetUserId(),
				DisplayName: in.GetDisplayName(),
				RoomID:      in.GetRoom(),
			})
			if err != nil {
				return translateError(err)
			}

			hasSession = true

			eventsCtx, cancel := context.WithCancel(ctx)
			eventsCancel = cancel
			eventsWG.Add(1)
			go s.forwardEvents(eventsCtx, &eventsWG, events, send, eventErr)

			if err := send(&chatv1.ServerEvent{
				Event: &chatv1.ServerEvent_Joined{
					Joined: &chatv1.JoinAck{
						UserId:         session.UserID,
						Room:           session.RoomID,
						WelcomeMessage: fmt.Sprintf(welcomeMessageFormat, session.DisplayName),
					},
				},
			}); err != nil {
				return err
			}

		case *chatv1.ClientEnvelope_Chat:
			if !hasSession {
				return status.Error(codes.FailedPrecondition, errMsgJoinRequired)
			}
			payload := msg.Chat
			if payload == nil {
				return status.Error(codes.InvalidArgument, errMsgChatPayloadRequired)
			}

			sentAt := time.Unix(payload.GetTimestampUtc(), 0).UTC()
			if payload.GetTimestampUtc() == zeroUnixTimestamp {
				sentAt = time.Now().UTC()
			}

			if err := s.chat.Broadcast(ctx, domain.Message{
				UserID:      payload.GetUserId(),
				DisplayName: "",
				RoomID:      payload.GetRoom(),
				Content:     payload.GetContent(),
				SentAt:      sentAt,
			}); err != nil {
				return translateError(err)
			}

		case *chatv1.ClientEnvelope_Leave:
			if !hasSession {
				return status.Error(codes.FailedPrecondition, errMsgNoActiveSession)
			}
			leave := msg.Leave
			if leave == nil {
				return status.Error(codes.InvalidArgument, errMsgLeavePayloadReq)
			}
			if err := s.chat.Leave(ctx, leave.GetRoom(), leave.GetUserId()); err != nil {
				return translateError(err)
			}
			hasSession = false
			return nil

		default:
			return status.Error(codes.InvalidArgument, errMsgInvalidPayload)
		}
	}
}

func (s *Server) forwardEvents(ctx context.Context, wg *sync.WaitGroup, events <-chan domain.Event, send func(*chatv1.ServerEvent) error, errCh chan<- error) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			if evt := domainEventToProto(ev); evt != nil {
				if err := send(evt); err != nil {
					select {
					case errCh <- err:
					default:
					}
					return
				}
			}
		}
	}
}

func domainEventToProto(ev domain.Event) *chatv1.ServerEvent {
	switch ev.Type {
	case domain.EventMessage:
		return &chatv1.ServerEvent{
			Event: &chatv1.ServerEvent_Broadcast{
				Broadcast: &chatv1.ChatPayload{
					UserId:       ev.UserID,
					Room:         ev.RoomID,
					Content:      ev.Content,
					TimestampUtc: ev.Timestamp.UnixMilli(),
				},
			},
		}
	case domain.EventUserJoined:
		return &chatv1.ServerEvent{
			Event: &chatv1.ServerEvent_Notice{
				Notice: &chatv1.ServerNotice{
					Type:    chatv1.ServerNotice_TYPE_USER_JOINED,
					Message: fmt.Sprintf(noticeJoinedFormat, ev.DisplayName),
					UserId:  ev.UserID,
					Room:    ev.RoomID,
				},
			},
		}
	case domain.EventUserLeft:
		return &chatv1.ServerEvent{
			Event: &chatv1.ServerEvent_Notice{
				Notice: &chatv1.ServerNotice{
					Type:    chatv1.ServerNotice_TYPE_USER_LEFT,
					Message: fmt.Sprintf(noticeLeftFormat, ev.DisplayName),
					UserId:  ev.UserID,
					Room:    ev.RoomID,
				},
			},
		}
	case domain.EventSystem:
		return &chatv1.ServerEvent{
			Event: &chatv1.ServerEvent_Notice{
				Notice: &chatv1.ServerNotice{
					Type:    chatv1.ServerNotice_TYPE_GENERIC,
					Message: ev.Content,
					UserId:  ev.UserID,
					Room:    ev.RoomID,
				},
			},
		}
	default:
		return nil
	}
}

func translateError(err error) error {
	switch {
	case errors.Is(err, usecase.ErrEmptyFields):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, usecase.ErrEmptyMessage):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, usecase.ErrAlreadyJoined):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, usecase.ErrRoomNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, usecase.ErrUserNotInRoom):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
