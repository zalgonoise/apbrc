package modifiers

import (
	"context"

	"github.com/zalgonoise/apbrc/log"
	"golang.org/x/exp/slog"
)

func AttributeWithLogs(a Attribute, logger log.Logger) Attribute {
	return attributeWithLogs{a, logger}
}

type attributeWithLogs struct {
	a      Attribute
	logger log.Logger
}

func (a attributeWithLogs) Match(ctx context.Context, line []byte) (ok bool, match string) {
	ok, match = a.a.Match(ctx, line)
	if ok {
		a.logger.InfoContext(ctx,
			"matched config attribute key",
			slog.String("key", match),
			slog.String("original_value", string(line)),
		)
	}

	return ok, match
}

func (a attributeWithLogs) Value(ctx context.Context, key string) (data []byte, value any) {
	data, value = a.a.Value(ctx, key)
	a.logger.InfoContext(ctx,
		"modified config attribute value",
		slog.String("key", key),
		slog.Any("new_value", value),
	)

	return data, value
}
