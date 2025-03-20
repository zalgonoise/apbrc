package log

import (
	"log/slog"

	"github.com/zalgonoise/cfg"
)

func WithSpanID() cfg.Option[*SpanContextHandler] {
	return cfg.Register[*SpanContextHandler](func(config *SpanContextHandler) *SpanContextHandler {
		config.withSpanID = true

		return config
	})
}

func WithHandler(handler slog.Handler) cfg.Option[*SpanContextHandler] {
	if handler == nil {
		handler = defaultHandler()
	}

	return cfg.Register[*SpanContextHandler](func(config *SpanContextHandler) *SpanContextHandler {
		config.handler = handler

		return config
	})
}
