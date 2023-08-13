package modifiers

import (
	"context"

	"github.com/zalgonoise/apbrc/log"
	"golang.org/x/exp/slog"
)

type Applier interface {
	Apply(ctx context.Context, basePath string) error
}

func ModifierWithLogs(m Modifier, logger log.Logger) Applier {
	if logger == nil {
		return m
	}

	attrs := make([]Attribute, 0, len(m.Attributes))
	for i := range m.Attributes {
		attrs = append(attrs, AttributeWithLogs(m.Attributes[i], logger))
	}

	return modifierWithLogs{NewModifier(m.FilePath, attrs...), logger}
}

type modifierWithLogs struct {
	m      Modifier
	logger log.Logger
}

func (m modifierWithLogs) Apply(ctx context.Context, basePath string) error {
	path := basePath + m.m.FilePath

	m.logger.InfoContext(ctx,
		"applying modifiers to config file",
		slog.String("path", path),
		slog.Int("num_attributes", len(m.m.Attributes)),
	)

	err := m.m.Apply(ctx, basePath)
	switch err {
	case nil:
		m.logger.InfoContext(ctx,
			"overwritten configuration file successfully",
			slog.String("path", path),
		)
	default:
		m.logger.ErrorContext(ctx,
			"failed to overwrite configuration file",
			slog.String("error", err.Error()),
			slog.String("path", path),
		)
	}

	return err
}
