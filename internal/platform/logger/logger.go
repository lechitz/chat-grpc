// Package logger provides a thin abstraction over zap for contextual logging.
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger exposes the logging contract used across the service.
type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Sync() error
}

// ZapLogger is a zap-backed implementation of Logger.
type ZapLogger struct {
	*zap.SugaredLogger
}

// Debugw proxies to the underlying sugared logger.
func (l *ZapLogger) Debugw(msg string, kv ...interface{}) {
	l.SugaredLogger.Debugw(msg, kv...)
}

// Infow proxies to the underlying sugared logger.
func (l *ZapLogger) Infow(msg string, kv ...interface{}) {
	l.SugaredLogger.Infow(msg, kv...)
}

// Warnw proxies to the underlying sugared logger.
func (l *ZapLogger) Warnw(msg string, kv ...interface{}) {
	l.SugaredLogger.Warnw(msg, kv...)
}

// Errorw proxies to the underlying sugared logger.
func (l *ZapLogger) Errorw(msg string, kv ...interface{}) {
	l.SugaredLogger.Errorw(msg, kv...)
}

// Sync flushes buffered log entries.
func (l *ZapLogger) Sync() error {
	return l.SugaredLogger.Sync()
}

// New constructs a zap logger tuned for local development by default.
// When env is production it switches to the production encoder.
func New(env string) (*ZapLogger, func(), error) {
	var (
		cfg zap.Config
		err error
	)

	if env == envProductionValue {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = productionTimeKey
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("build zap logger: %w", err)
	}

	sugar := logger.Sugar()
	cleanup := func() {
		_ = sugar.Sync()
	}

	return &ZapLogger{SugaredLogger: sugar}, cleanup, nil
}
