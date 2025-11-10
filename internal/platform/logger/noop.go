// Package logger provides a thin abstraction over zap for contextual logging.
package logger

import (
	"context"

	portslogger "github.com/lechitz/chat-grpc/internal/platform/ports/logger"
)

// NoopLogger is a reusable no-op implementation of Logger and the
// ports logger ContextLogger intended for tests and places where logging
// should be silenced.
//
// Use in tests as: logs := logger.NoopLogger{}
// It intentionally does nothing and is safe to use in concurrent tests.
type NoopLogger struct{}

func (NoopLogger) Infof(string, ...any)  {}
func (NoopLogger) Errorf(string, ...any) {}
func (NoopLogger) Debugf(string, ...any) {}
func (NoopLogger) Warnf(string, ...any)  {}

func (NoopLogger) Infow(string, ...any)  {}
func (NoopLogger) Errorw(string, ...any) {}
func (NoopLogger) Debugw(string, ...any) {}
func (NoopLogger) Warnw(string, ...any)  {}

func (NoopLogger) InfowCtx(context.Context, string, ...any)  {}
func (NoopLogger) ErrorwCtx(context.Context, string, ...any) {}
func (NoopLogger) WarnwCtx(context.Context, string, ...any)  {}
func (NoopLogger) DebugwCtx(context.Context, string, ...any) {}

func (NoopLogger) Sync() error { return nil }

var _ Logger = (*NoopLogger)(nil)
var _ portslogger.ContextLogger = (*NoopLogger)(nil)
