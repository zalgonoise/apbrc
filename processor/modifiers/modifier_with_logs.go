package modifiers

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/zalgonoise/apbrc/monitoring"
)

type Applier interface {
	Apply(ctx context.Context, basePath string) error
}

func ModifierWithLogs(m Modifier, logger monitoring.Logger) Applier {
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
	logger monitoring.Logger
}

func (m modifierWithLogs) Apply(ctx context.Context, basePath string) error {
	path := m.m.FilePath + basePath

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
