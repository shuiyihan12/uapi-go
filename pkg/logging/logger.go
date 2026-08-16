// Package logging provides a zap-based enterprise logging facade with
// structured fields, raw-text output and trace-ID propagation.
package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

// noopLogger discards everything — the zero-dependency default for library
// consumers that bring their own logging stack.
type noopLogger struct{}

// Noop returns a Logger that discards all output. It is the recommended
// default for SDK consumers that do not want uapi-go to log on its own.
func Noop() Logger { return noopLogger{} }

func (noopLogger) Debug(string, ...zap.Field) {}
func (noopLogger) Info(string, ...zap.Field)  {}
func (noopLogger) Error(string, ...zap.Field) {}
func (noopLogger) Warn(string, ...zap.Field)  {}
func (noopLogger) Raw(string)                 {}
func (n noopLogger) With(...zap.Field) Logger { return n }

// WithContext keeps returning the noop logger.
func (n noopLogger) WithContext(context.Context) Logger { return n }

// stdLogAdapter adapts the standard library *log.Logger to the Logger
// interface for consumers without zap.
type stdLogAdapter struct {
	l *log.Logger
}

// NewStdLogAdapter wraps a standard library *log.Logger into a Logger.
// Structured fields are flattened onto the line as key=value pairs.
func NewStdLogAdapter(l *log.Logger) Logger { return stdLogAdapter{l: l} }

func (a stdLogAdapter) log(level, msg string, fields []zap.Field) {
	var b strings.Builder
	b.WriteString(level)
	b.WriteByte(' ')
	b.WriteString(msg)
	for _, f := range fields {
		b.WriteByte(' ')
		b.WriteString(f.Key)
		b.WriteByte('=')
		val := f.String
		if val == "" {
			switch {
			case f.Interface != nil:
				val = fmt.Sprintf("%v", f.Interface)
			case f.Integer != 0:
				val = strconv.FormatInt(f.Integer, 10)
			}
		}
		b.WriteString(val)
	}
	a.l.Print(b.String())
}

func (a stdLogAdapter) Debug(msg string, fields ...zap.Field) { a.log("DEBUG", msg, fields) }
func (a stdLogAdapter) Info(msg string, fields ...zap.Field)  { a.log("INFO", msg, fields) }
func (a stdLogAdapter) Error(msg string, fields ...zap.Field) { a.log("ERROR", msg, fields) }
func (a stdLogAdapter) Warn(msg string, fields ...zap.Field)  { a.log("WARN", msg, fields) }

// Raw writes the payload verbatim.
func (a stdLogAdapter) Raw(msg string) { a.l.Print(msg) }

func (a stdLogAdapter) With(...zap.Field) Logger { return a }

// WithContext keeps returning the same adapter; trace fields are dropped.
func (a stdLogAdapter) WithContext(context.Context) Logger { return a }
