package monitoring

import (
	"context"
)

type Logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

func NoopLogger() noopLogger {
	return noopLogger{}
}

type noopLogger struct{}

func (l noopLogger) DebugContext(ctx context.Context, msg string, args ...any) {}
func (l noopLogger) InfoContext(ctx context.Context, msg string, args ...any)  {}
func (l noopLogger) WarnContext(ctx context.Context, msg string, args ...any)  {}
func (l noopLogger) ErrorContext(ctx context.Context, msg string, args ...any) {}

type Observability interface {
	Logger
}
