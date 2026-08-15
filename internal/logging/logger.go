// Package logging provides a zap-based enterprise logging facade with
// structured fields, raw-text output and trace-ID propagation.
package logging

import (
	"context"
	"fmt"
	"os"

	"github.com/shuiyihan12/uapi-go/pkg/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the enterprise logging interface.
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	// Raw writes one line of raw text (no JSON encoding), suitable for
	// printing large XML payloads.
	Raw(msg string)
	With(fields ...zap.Field) Logger
	WithContext(ctx context.Context) Logger
}

// zapLogger is the zap-backed Logger implementation.
type zapLogger struct {
	logger *zap.Logger
}

// NewLogger creates a new logger.
func NewLogger(level string, isDevelopment bool) (Logger, error) {
	var config zap.Config

	if isDevelopment {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.Encoding = "json"
	}

	// Apply the requested log level.
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := config.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: logger}, nil
}

// NewDefaultLogger creates the default logger.
func NewDefaultLogger() Logger {
	isDev := os.Getenv("ENVIRONMENT") != "production"
	logger, _ := NewLogger("info", isDev)
	return logger
}

// Info writes one structured log entry at Info level.
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Error writes one structured log entry at Error level.
func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Debug writes one structured log entry at Debug level.
func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Warn writes one structured log entry at Warn level.
func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Raw writes one raw-text log line, bypassing JSON encoding; suitable for
// printing large XML payloads.
func (l *zapLogger) Raw(msg string) {
	if !l.logger.Core().Enabled(zapcore.InfoLevel) {
		return
	}
	fmt.Fprintln(os.Stderr, msg)
}

// With returns a new Logger carrying additional structured fields.
func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

// WithContext returns a Logger whose log fields include the trace ID from
// the context, enabling end-to-end correlation.
func (l *zapLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}
	// Surface the global trace ID as a log field for full request-response
	// correlation.
	if id := trace.ID(ctx); id != "" {
		return &zapLogger{logger: l.logger.With(zap.String("trace_id", id))}
	}
	return l
}
